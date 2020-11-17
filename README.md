# Prometheus exporter for DNT roomlogg-pro

This is a small application to read temperature and humidity sensor data from a connected [dnt RoomLogg Pro](https://www.dnt.de/Produkte/Raumklimastation-RoomLogg-PRO/) device and provide it as prometheus metrics exporter.
The logic is based on the well documented reverse-engineering done by the [Raumklima](https://github.com/juergen-rocks/raumklima) tool.

## Usage

The exporter tool provides two different operation modes: single query and server mode.

The single query mode can be enabled using `-query`. It will print the current sensor data and exit afterwards.

In order to start a prometheus exporter, you can use the parameter `-server`. By default the server runs on port 8080, but you can change it using the parameter `-port xxx` where xxx is a valid port number. Afterwards you can query the metrics using e.g. `curl http://localhost:8080/metrics`.

If you encounter a message that HID is not supported, make sure that the application was compiled correctly for your platform (see information regarding cross compilation below). If opening the device fails, make sure that you have the correct access rights to read from USB devices.  

## Setup

### USB

The application needs USB access on a HID device. Therefore you either have to run the compiled application as root or create a udev rule to make the device accessible.

To make it accessible via udev, create a file called `/etc/udev/rules.d/70-roomlogg-pro.rules` with the following content:

```
SUBSYSTEM=="usb", ATTRS{idVendor}=="0483", ATTRS{idProduct}=="5750", MODE="0666"
```

Save the file and run `udevadm control --reload-rules`. You can find more information [here](https://askubuntu.com/questions/978552/how-do-i-make-libusb-work-as-non-root).

### Cross compilation

In order to run the application on an ARM device like the Raspberry Pi, you need to cross-compile the application for that specific architecture.
Since the go hid library is compiled from native C code, you need to install a gcc cross-compile toolchain in order to build that specific library. This is also described in the [library's repository](https://github.com/karalabe/hid).

I personally use [xgo](https://github.com/karalabe/xgo) to compile the application for a Raspberry Pi Zero. After installing xgo, you can create an ARM executable easily:

```
xgo -out build/roomlogg --targets=linux/arm -x github.com/jhendess/roomlogg-go/cmd/
```
