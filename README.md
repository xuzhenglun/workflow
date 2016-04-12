# workflow
Simple workflow engine implement via Golang and Lua

一个用Golang配合Lua脚本写的工作流引擎

Install
----
```bash
go get github.com/xuzhenglun/workflow
go run main.go
go run front/main.go
```

Config
---

1. 单个逻辑动作写成一个lua脚本，格式如下，保存在scripts下：

    文件名:start.lua
    ```Lua
    start = {  					--对象名，需与文件名同,start为固定事务起点
    	father = null，			--父事务名（初始事务不需要定义）
        name = "start",			--事务名，需与文件名同
        typ = "POST",				--事务处理方式，POST/GET
        groups = "start",			--用户组
        needArgs = "name address idcard",	 --被修改参数，任意特殊字符分隔
        handler = function(req)	             --handler函数，一个形参，为请求
            log("this is starter")	         --调用Api
            return AddRow(start,req)		 --调用Api，返回结果
        end
    }
    ```
2. 证书保存在工作文件夹下的keys.json中：

    ```json
    {
        "Key":
            {
                "Prikey":"\n-----BEGIN RSA PRIVATE KEY-----\nMIICXgIBAAKBgQDzjHVNKSeiKibvblbNVCFANsGWCFKhNnXEA5kkYiaDxwDCXKQo\n3P9hVwKdWVOa9IgOmsjuJithiCgQnz204ReQ1soPXor5/h8HGMxVRgBMEcb0Hjgg\n4yURTBjaStm1G6DIMywGEuNtpJ+/XZeeQunlceAf2ZYA60BDcKDdadv39wIDAQAB\nAoGBAJpjuuSMJ0TEpeP4NWx6XY3AaF28ruzlgigdA9Ktqa6104Rih+ojlnzVVKH8\nw2QiibGAa8mURsHQN126JLrqSi3ZYtMWfmUz6bbw0ZtPKhoRgkxGbTGGx9LngxSA\nEBz+A6HVSFtw091C/pGgM7zybMNjsr7fE631Xn8Jgyca0+zpAkEA/1j79wpxFvR6\nROPUCwJ7zLBa/+r5dQeh8zHVyShR09xF0/ok6CXOykI1j7KF3un0bfghlsZCZISx\nl8iut+iExQJBAPQrwbG+tg5hs2rDRGBBtHgoF+WbsFyy85TiIyF+iFTlHQHxuPvE\nFXp3ov1ggbqYlmZ6kjlOaoaU/j4cBz8EbYsCQQC8AEwjKzDwcbfEGOn54S49Gmsl\nmV9pZuE6KSr6HBGDJt7Sn42kzpFeITlGP26JHT+158bzN62STJBk7ICXLz7xAkBB\n6FW+PrYxp5mgZdjdCp9GF7xrk9zFCODK/UdyUQ9HqxhrX+4It2L8zbJHJneeAHYI\nb2ls4ofKkAkYhsRF9FIFAkEAlC2U1NpEcIbj5uGFrxS13eKfmoZ5UdwmeSzH/YnV\n2LIc2EJFQt6A6EoZPUtRJw0WYB49NmLIQ4uU1uN/hC+5Rg==\n-----END RSA PRIVATE KEY-----\n",
                "Pubkey":"\n-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDzjHVNKSeiKibvblbNVCFANsGW\nCFKhNnXEA5kkYiaDxwDCXKQo3P9hVwKdWVOa9IgOmsjuJithiCgQnz204ReQ1soP\nXor5/h8HGMxVRgBMEcb0Hjgg4yURTBjaStm1G6DIMywGEuNtpJ+/XZeeQunlceAf\n2ZYA60BDcKDdadv39wIDAQAB\n-----END PUBLIC KEY-----\n"
            },
        "Method":"RSA"
    }
    ```
    
---
