package opts

type Opts struct {
	File string `short:"f" long:"file" description:"Game file to load" required:"true"`
	//Should change this to a string so as to handle hex values
	//Blue 23455
	BgColour uint32 `short:"b" long:"bg-colour" description:"Background Colour" required:"false"`
	XSize    int    `long:"x-size" description:"Effective pixel size in pixels"`
	YSize    int    `long:"y-size" description:"Effective pixel size in pixels"`
}
