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
$ git clone https://github.com/ochinchina/minimalDockerImage.git
$ cd minimalDockerImage
$ export GOPATH=`pwd`
$ go install mdi
$ export PATH=$GOPATH/bin:$PATH
```
##run the tool to make package for you

###make a tarball file

Create a tarball to include /usr/bin/ls, /usr/bin/pwd file and their dependencies, execute following command:

```shell
$mdi -i /usr/bin/ls,/usr/bin/pwd -o test.tar create_tar
```

If too many files you want to include in the tarball file, you can put them into a file, for example put to input.txt

```shell
$ cat input.txt
/usr/bin/ls
/usr/bin/pwd
```
then use flag "-f file_name" to execute the command:

```shell
$ mdi -f input.txt -o test.tar create_tar
```

Except for the executable binary file, if you want to pack a directory or some non-executable binary file, this tool can also accept the directory file or non-executable file as input. If input a directory, this tool will try to find all the executable binary files and the .so files under the directory recursively, then find the dependency files for you.

For example, if you want to pack directory "/usr/bin" to the tarball, you can simply run:

```shell
$ mdi -i /usr/bin -o test.tar create_tar
```

if some .so files or executable binary files under /usr/bin depend on .so files in other directory, the output test.tar file will contains them also.

###make docker image

If docker is started in the host running mdi tool, you can make a docker image directly, the optional flag "-i" and "-f" is same as the create_tar command.

For example, if you want to make a docker image named "haproxy" to include the /usr/bin/haproxy, you can run:

```shell
$ mdi -i /usr/bin/haproxy -o haproxy create_image
```

###print all the files

If you just want to see what files will be contained in the docker image, you can run the tool like:

```shell
$ mdi -i /usr/bin/haproxy list_files
```

the above command will list all the files will be included in the docker image

## License

The MIT License (MIT)

Copyright (c) <year> <copyright holders>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.


