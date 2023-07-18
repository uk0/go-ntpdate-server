# go-ntpdate-server
simple ntpdate server



## quick start 

* `ntpdate -d -u 127.0.0.1`

```bash
root@ubuntu-test-58:~# ntpdate -d -u 127.0.0.1
18 Jul 04:02:51 ntpdate[1682912]: ntpdate 4.2.8p15@1.3728-o Wed Feb 16 17:13:02 UTC 2022 (1)
Looking for host 127.0.0.1 and service ntp
127.0.0.1 reversed to localhost
host found : localhost
transmit(127.0.0.1)
receive(127.0.0.1)
transmit(127.0.0.1)
receive(127.0.0.1)
transmit(127.0.0.1)
receive(127.0.0.1)
transmit(127.0.0.1)
receive(127.0.0.1)

server 127.0.0.1, port 123
stratum 1, precision -6, leap 00, trust 000
refid [], root delay 0.000000, root dispersion 0.000000
reference time:      e8608d71.8cac5b44  Tue, Jul 18 2023  4:02:57.549
originate timestamp: e8608d71.8cac6964  Tue, Jul 18 2023  4:02:57.549
transmit timestamp:  e8608d71.9702e3d3  Tue, Jul 18 2023  4:02:57.589
filter delay:  0.04166    0.04149    0.04150    0.04149
               ----       ----       ----       ----
filter offset: -0.040514  -0.040481  -0.040477  -0.040505
               ----       ----       ----       ----
delay 0.04149, dispersion 0.00002, offset -0.040481

18 Jul 04:02:57 ntpdate[1682912]: adjust time server 127.0.0.1 offset -0.040481 sec
```