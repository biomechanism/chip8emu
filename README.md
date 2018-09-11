# Chip-8 Emulator/VM in Go

![Pong](/doc/images/chip8emu-1.png)

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


## Tehnical Details

The initial implementation of the drawing functionality was a bit unwieldy and convoluted, I’ve
simplified and cleaned things up considerably since then. Partially it comes for the challenge of just
sitting down and figuring the algorithms out myself, and it worked fine, but it was a bit overly
complex for something so simple. Initially my objective was to only use the actual 4096? bytes of
memory, I have since scrapped that and seperated the video memory, and defined it as a
multidimensional array for even easier access. I’m also using 1 byte per pixel, as opposed to 1 byte 
per 8 pixels in the original bit packed memory, it makes life a bit (no pun intended) easier.

### Clock Rate

There doesn’t appear to be too much available on the most appropriate clock frequency, I’ve generally seen it 
mentioned in other material of a rate around 500Hz – 540Hz. This emulator has it clocked at 500Hz, theoretically, and I say theoretically because I’ve chosen to use sdl.Delay() at two milliseconds. However with this approach the timing is really at the mercy and vagaries of the underlying OS scheduler. I have found it quite adequate however.
  
I had also considered allowing a tunable clock frequency, but ultimately decided against it, as it
turned out to be a trade-off between a more adjustable emulation speed or lower CPU utilisation.
I found using Go’s time.Sleep() and channel based approaches to be pretty equivalent, and quite
flexible as the allow much more granularity with regards to timing. The downside of this however
was roughly 10 – 15%+ CPU utilisation when compared with sdl.Delay().
