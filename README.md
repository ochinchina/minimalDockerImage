# minimalDockerImage

How do you create your docker image? For most people, they download a base docker image, 
such as ubuntu and then add their application upon the base image. This works good but
the image has little bit big. Sometimes you don't really need a whole base image, you
just simply want to make your docker application run correctly.

This tool will help you make minimal docker image. Before using this tool, you should make
sure your application can run in one linux environment and then input your application absolute
path to the tool. This tool will package the application and all its dependency .so files to
a tarball or make a docker image directly if you have installed your docker in your running
environment.

##compile the tool

```shell
$ export GOPATH=<your go path>
$ go get minimalDockerImage
```

##run the tool to make package for you

###make a tarball file
###print all the files
###make docker image
