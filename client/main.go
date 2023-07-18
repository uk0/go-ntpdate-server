package main

import (
	"encoding/binary"
	"fmt"
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

// Convert NTP timestamp format to time.Time
func fromNTPTime(ntpTime [8]byte) time.Time {
	seconds := binary.BigEndian.Uint32(ntpTime[:4])
	fraction := binary.BigEndian.Uint32(ntpTime[4:])
	nanoseconds := time.Duration(seconds)*time.Second + time.Duration(fraction)*time.Second/0x100000000
	return ntpEpoch.Add(nanoseconds)
}

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:123")
	if err != nil {
		log.Fatalf("Error resolving address: %v", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("Error dialing UDP: %v", err)
	}
	defer conn.Close()

	packet := &ntpPacket{
		Settings: 0x23, // LI = 0, VN = 4, Mode = 3 (client)
	}

	b := make([]byte, 48)
	b[0] = packet.Settings

	_, err = conn.Write(b)
	if err != nil {
		log.Fatalf("Error writing to UDP: %v", err)
	}

	_, err = conn.Read(b)
	if err != nil {
		log.Fatalf("Error reading from UDP: %v", err)
	}

	packet = &ntpPacket{}
	packet.Settings = b[0]
	packet.Stratum = b[1]
	packet.Poll = b[2]
	packet.Precision = b[3]
	packet.RootDelay = binary.BigEndian.Uint32(b[4:8])
	packet.RootDispersion = binary.BigEndian.Uint32(b[8:12])
	packet.ReferenceID = binary.BigEndian.Uint32(b[12:16])
	copy(packet.ReferenceTimestamp[:], b[16:24])
	copy(packet.OriginateTimestamp[:], b[24:32])
	copy(packet.ReceiveTimestamp[:], b[32:40])
	copy(packet.TransmitTimestamp[:], b[40:48])

	fmt.Printf("NTP server's time is: %v\n", fromNTPTime(packet.TransmitTimestamp))
	fmt.Printf("NTP server's time is: %v\n", fromNTPTime(packet.ReceiveTimestamp))
	fmt.Printf("NTP server's time is: %v\n", fromNTPTime(packet.OriginateTimestamp))
}
