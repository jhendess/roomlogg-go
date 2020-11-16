# Prometheus exporter for DNT roomlogg-pro

This is a tool to read temperature and humidity sensor data from a connected [Roomlogg PRO](https://www.reichelt.com/de/en/roomlogg-pro-room-climate-control-station-dnt-roomlogg-pro-p267987.html) device.

## Setup

### USB

The application needs USB access on a non-default HID device under linux. Therefore you either have to run the compiled application as root or create a udev rule to make the device accessible.

To make it accessible via udev, create a file called `/etc/udev/rules.d/70-roomlogg-pro.rules` with the following content:

```
SUBSYSTEM=="usb", ATTRS{idVendor}=="0483", ATTRS{idProduct}=="5750", MODE="0666"
```

Save the file and run `udevadm control --reload-rules`. You can find more information [here](https://askubuntu.com/questions/978552/how-do-i-make-libusb-work-as-non-root).

### Crosscompiling

In order to run the application on an ARM device like the Raspberry Pi, you need to cross-compile the application for that specific architecture.
Since the go hid library is compiled from native C code, you need to install a gcc cross-compile toolchain in order to build that specific library. This is also described in the [library's repository](https://github.com/karalabe/hid).

On Arch Linux you can install the `arm-linux-gnueabi-gcc` package from the AUR

Afterwards run the following command to create your ARM executable:


