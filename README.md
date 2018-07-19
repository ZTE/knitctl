# knitctl
A awesome command line tool project to manager knitter API
# How to use
   - knitctl help
   - knitctl version
   - knitctl config [options]
   - knitctl create [options]
   - knitctl create [(type [NAME])] [options]
   - knitctl get [(type [NAME])] [options]
   - knitctl set [(type [NAME])] [options]
   - knitctl delete [options]
   - knitctl delete [(type [NAME])] [options]

## help  
   help to use command line tool
   - knitctl help
  
## version
   version of knitter api
   - knitctl version
   
## config
config knitter sever ip and port
   - knitctl config -f config.yaml
   - knitctl config --knitterurl=localhost:9527
 
### config yaml example
a example config yaml file to config knitter server
```
knitterurl: localhost:9527
```

## create (file)
create resource(tenant/network/ipgroup) from a file or command

### from a yaml file or a folder
   - knitctl create -f resource.yaml
   - knitctl create -f -r /home/resource/
   
#### resource yaml example

```
kind: tenant
apiVersion: v1
metadata:
  name: my-tenant
spec:
---
kind: network
apiVersion: v1
metadata:
  name: my-nw
  tenant: my-tenant
spec:
  gateway: 10.144.209.1
  cidr: 10.144.209.0/24
  public: false
---
kind: ipgroup
apiVersion: v1
metadata:
  name: my-ipg
  tenant: my-tenant
  network: my-nw
spec:
  ips: 10.144.209.2
  
```


### create (command)

create tenant from command
   - knitctl create tenant my-tenant
create network from command
   - knitctl create network my-nw --tenant=my-tenant --public=false --gateway=10.144.209.1 --cidr=10.144.209.0/24
create ipgroup from command
   - knictl create ipgroup my-ipg --tenant=my-tenant --network=my-nw --ips=10.144.209.2,10.144.209.3
   
## get
list resource from command
   - knitctl get tenant
   - knitctl get tenant my-tenant
   - knitctl get network --tenant=my-tenant
   - knitctl get network my-nw --tenant=my-tenant
   - knitctl get ipgroup --tenant=my-tenant
   - knitctl get ipgroup my-ipg --tenant=my-tenant
   - knitctl get pod --tenant=my-tenant
   - knitctl get pod my-pod --tenant=my-tenant
   
## set
update resource from command
   - knitctl set tenant my-tenant --quota=3
   - knitctl set ipgroup my-ipg --tenant=my-tenant --ips=10.144.209.2