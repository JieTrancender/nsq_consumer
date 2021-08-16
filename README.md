**This is a tool which consumes nsqd messages to other consumer.**

### Quick Start
1. Build nsqd environment.
2. Build etcd environment.
3. Put following configure value to etcd's key `/config/nsq_consumer/default`.
~~~
{
    "consumer-name":"nsq-consumer",
    "nsq-consumer": {
        "lookupd-http-addresses":["http://127.0.0.1:4161"],
        "topics":[
            "dev_test",
            "dev_test_2"
        ],
        "consumer-type":"nsq"
    },
    "output": {
        "console": {
            "enabled": true
        },
        "nsqd": {
            "nsqd": "127.0.0.1:4150",
            "topic": "dev_test_dup",
            "enabled": false,
            "enabled_topic": true
        },
        "elasticsearch": {
            "enabled": false,
            "addrs": ["http://127.0.0.1:9200"],
            "username": "root",
            "password": "123456"
        }
    },
    "logging": {
        "level": 0,
        "to_stderr": true
    }
}
~~~
4. Run this project
~~~
make clean && ./build/nsq_to_consumer --etcd-endpoints 127.0.0.1:2379 --etcd-username root --etcd-password 123456 --etcd-path /nsq_consumer/default
~~~
### Output list
1. console
2. nsqd
3. elasticsearch
4. file[todo]
5. http[todo]
6. mysql[todo]

### Getting Help
If you need help or hit an issue, you can make a issue, we will deal it as soon as posibile.

### Contributing
We'd love working with you! You can do any thing if it's helpful, such as adding document, adding more consumer and so on.

### Note
This project's code makes in-depth reference to [beats](https://github.com/elastic/beats).
