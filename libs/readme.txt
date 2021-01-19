glibc:

        解压

        tar -zxvf  glibc-2.18.tar.gz

        创建编译目录

        cd glibc-2.28 

        mkdir build

        编译、安装

        cd build/

        ../configure --prefix=/usr --disable-profile --enable-add-ons --with-headers=/usr/include --with-binutils=/usr/bin

        make -j 8

        make install

        strings /lib64/libc.so.6 | grep GLIBC

libstdc++:
        把libstdc++.so.6.0.25拷贝到/usr/lib64目录下。

        　　cp libstdc++.so.6.0.25 /usr/lib64/

        删除原来的libstdc++.so.6符号连接。

        　　rm -rf libstdc++.so.6

        新建新符号连接。

        　　ln -s libstdc++.so.6.0.25 libstdc++.so.6

