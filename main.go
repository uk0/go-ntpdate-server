package main

import (
	"encoding/binary"
	"golang.org/x/net/ipv4"
	"log"
	"net"
	"time"
)

// NTP epoch (start of time for network time protocol) is 01/01/1900.
var ntpEpoch = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)

// NTP packet structure
type ntpPacket struct {
	Settings           uint8 // LeapIndicator, VersionNumber, Mode
	Stratum            uint8
	Poll               uint8
	Precision          uint8
	RootDelay          uint32
	RootDispersion     uint32
	ReferenceID        uint32
	ReferenceTimestamp [8]byte
	OriginateTimestamp [8]byte
	ReceiveTimestamp   [8]byte
	TransmitTimestamp  [8]byte
}

// Convert time.Time to NTP timestamp format
func toNTPTime(t time.Time) (ntpTime [8]byte) {
	seconds := uint32(t.Sub(ntpEpoch).Seconds())
	fraction := uint32((t.Sub(ntpEpoch) % time.Second) * (1 << 32 / time.Second))
	binary.BigEndian.PutUint32(ntpTime[:4], seconds)
	binary.BigEndian.PutUint32(ntpTime[4:], fraction)
	return ntpTime
}

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":123")
	if err != nil {
		log.Fatalf("Error resolving address: %v", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("Error listening on UDP: %v", err)
	}
	defer conn.Close()

	p := ipv4.NewPacketConn(conn)
	defer p.Close()

	b := make([]byte, 48)
	for {
		_, _, remoteAddr, err := p.ReadFrom(b)
		if err != nil {
			log.Fatalf("Error reading from UDP: %v", err)
		}

		var packet *ntpPacket
		// Parse the client's request
		clientRequest := ntpPacket{}
		clientRequest.Settings = b[0]
		clientRequest.Stratum = b[1]
		clientRequest.Poll = b[2]
		clientRequest.Precision = b[3]
		clientRequest.RootDelay = binary.BigEndian.Uint32(b[4:8])
		clientRequest.RootDispersion = binary.BigEndian.Uint32(b[8:12])
		clientRequest.ReferenceID = binary.BigEndian.Uint32(b[12:16])
		copy(clientRequest.ReferenceTimestamp[:], b[16:24])
		copy(clientRequest.OriginateTimestamp[:], b[24:32])
		copy(clientRequest.ReceiveTimestamp[:], b[32:40])
		copy(clientRequest.TransmitTimestamp[:], b[40:48])

		mode := b[0] & 0x7
		switch mode {
		case 3: // Client mode

			// Create the server's response
			packet = &ntpPacket{
				Settings:           0x24,                            // LI = 0, VN = 4, Mode = 4 (server)
				Stratum:            1,                               // Primary reference (e.g., an atomic clock)
				Poll:               clientRequest.Poll,              // Same as client's request
				Precision:          clientRequest.Precision,         // Same as client's request
				RootDelay:          0,                               // Server is directly connected to a reference clock
				RootDispersion:     0,                               // Server is directly connected to a reference clock
				ReferenceID:        0,                               // Server is directly connected to a reference clock
				ReferenceTimestamp: toNTPTime(time.Now()),           // Current time
				OriginateTimestamp: clientRequest.TransmitTimestamp, // Same as client's transmit timestamp
				ReceiveTimestamp:   toNTPTime(time.Now()),           // Server receive time
				TransmitTimestamp:  toNTPTime(time.Now()),           // Server transmit time
			}

			log.Println("mode 3 ", mode)
			copy(b[32:40], packet.TransmitTimestamp[:])

		case 4: // Server mode
			log.Println("Received a packet in server mode. Ignoring it.")
		default:
			log.Printf("Received a packet with unknown mode: %v. Ignoring it.", mode)
		}

		// Convert the response to bytes
		b[0] = packet.Settings
		b[1] = packet.Stratum
		b[2] = packet.Poll
		b[3] = packet.Precision
		binary.BigEndian.PutUint32(b[4:8], packet.RootDelay)
		binary.BigEndian.PutUint32(b[8:12], packet.RootDispersion)
		binary.BigEndian.PutUint32(b[12:16], packet.ReferenceID)
		copy(b[16:24], packet.ReferenceTimestamp[:])
		copy(b[24:32], packet.OriginateTimestamp[:])
		copy(b[32:40], packet.ReceiveTimestamp[:])
		copy(b[40:48], packet.TransmitTimestamp[:])
		log.Println("send data to client ....")
		_, err = conn.WriteTo(b, remoteAddr)
		if err != nil {
			log.Fatalf("Error writing to UDP: %v", err)
		}
	}
}
