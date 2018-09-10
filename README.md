# Chip-8 Emulator/VM in Go

So the emulator is finally finished! Well, functional, I still need to make some fixes around the audio beep. 
It has taken a while, mostly because I’ve only done little bits now and then, and leaving it for many months when encountering particularly annyoying bugs. One in particular drove
me mad. It appears somewhere along the line the implementation of Chip 8 VMs has drified from
the original spec. regarding the shift right and left operations.

The original specs for those are vx = vy >>1 and vx = vy << 1, however all the recent VMs and
programs appear to operate on the logic of vx = vx >> 1 and vx = vx <<1.
That caused a bit of pain as many programs worked while particular problems manifested under
others, namly Invaders and Blinky.

## Downloading and Building
You will need git and Go installed

````bash
git clone <repo>
cd <directory>
go build
./<program>
````

## Supported parameters
```bash
-f <ROM to load> Specify the game file to load.
-b <value> Specify a decimal value to change the background colour
--x-size <value> The value to use is the emulated pixel size, in real pixels. Defaults to 10.
--y-size <value> The value to use is the emulated pixel size, in real pixels. Defaults to 10.

e.g.
<program> -b 34354 –x-size 15 –y-size 15
```

## Key Mapping
<table border="0">
<tr>
<td>

<table>
<tr> Keyboard
  <td> 1 </td>
  <td> 2 </td>
  <td> 3 </td>
  <td> 4 </td>
</tr>
<tr>  
  <td> Q </td>
  <td> W </td>
  <td> E </td>
  <td> R </td>
</tr>
<tr>
  <td> A </td>
  <td> S </td>
  <td> D </td>
  <td> F </td>
</tr>
<tr>
  <td> Z </td>
  <td> X </td>
  <td> C </td>
  <td> V </td>
</tr>

</table>
</td>

<td>
&nbsp
</td>

<td>
<table>
<tr> Chip-8 Keypad
  <td> 1 </td>
  <td> 2 </td>
  <td> 3 </td>
  <td> C </td>
</tr>
<tr>  
  <td> 4 </td>
  <td> 5 </td>
  <td> 6 </td>
  <td> D </td>
</tr>
<tr>
  <td> 7 </td>
  <td> 8 </td>
  <td> 9 </td>
  <td> E </td>
</tr>
<tr>
  <td> A </td>
  <td> 0 </td>
  <td> B </td>
  <td> F </td>
</tr>


</table>
</td>
</tr>
</table>

