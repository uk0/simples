---
layout: post
title: Python流量加速
categories: python
description: Python流量加速
keywords: linux,python
---


一个用Python做网卡流量加速

#  先看代码
date：2017-10-11
author：zhangjianxin

[TOC]


```python

#!/usr/bin/env py262

#导入方法模块
import os
import sys
import re

from xml.dom import minidom
from lxml import etree
from StringIO import StringIO
import time

file_eth1 = '/etc/sysconfig/network-scripts/ifcfg-eth1'
file_eth0 = '/etc/sysconfig/network-scripts/ifcfg-eth0'
file_sysctl = '/etc/sysctl_test.conf'
file_vd = '/etc/init.d/varnishd'
file_vf = '/usr/local/varnish_214/etc/varnish_01.conf'
filexml= '/speedy_data_transfers_server/speedy_data_transfers_server/WEB-INF/classes/config/config.xml'
nw = {}
eth = {}


def daemonize (stdin='/dev/null', stdout='/dev/null', stderr='/dev/null'):
    ''' Fork the current process as a daemon, redirecting standard file
        descriptors (by default, redirects them to /dev/null).
    '''
    # Perform first fork.
    try:
        pid = os.fork( )
        if pid > 0:
            sys.exit(0) # Exit first parent.
    except OSError, e:
        sys.stderr.write("fork #1 failed: (%d) %s\n" % (e.errno, e.strerror))
        sys.exit(1)
    # Decouple from parent environment.
    os.chdir("/")
    os.umask(0)
    os.setsid( )
    # Perform second fork.

    try:
        pid = os.fork( )
        if pid > 0:
            sys.exit(0) # Exit second parent.
    except OSError, e:
        sys.stderr.write("fork #2 failed: (%d) %s\n" % (e.errno, e.strerror))
        sys.exit(1)
    # The process is now daemonized, redirect standard file descriptors.
    for f in sys.stdout, sys.stderr: f.flush( )
    si = file(stdin, 'r')
    so = file(stdout, 'a+')
    se = file(stderr, 'a+', 0)
    os.dup2(si.fileno( ), sys.stdin.fileno( ))
    os.dup2(so.fileno( ), sys.stdout.fileno( ))
    os.dup2(se.fileno( ), sys.stderr.fileno( ))



def parse_xml(fn,r_t,*tg):
    f_x = open(fn,'r')
    f_x.seek(0,0)
    tree = etree.parse(fn)
    root = tree.getroot()
    #print root
    #print type(root)
    NS = '{http://accelerator.com/config}'

    r_t_element = root.findall( NS  + r_t)
    #print r_t_element
    #print type(r_t_element)
    if len(tg) >= 2:
        #if tg[-2] is None:
           #print tg[-2]
        s_t_element = tree.findall("//" + NS  + tg[-2])
        tg1_element = s_t_element[0].findall(NS + tg[-1])
        b = []
        for i in tg1_element:
            b.append(i.text)
        return b
    if len(tg) < 2:
        tg1_element = r_t_element[0].findall(NS + tg[-1])
        b = []
        for i in tg1_element:
            b.append(i.text)
        return b

def file_replace(filename,srcStr,dstStr):
    f = open(filename,'r')
    sc = f.read()
    stra = re.compile(srcStr)
    sub = re.sub(stra,dstStr,sc,0)
    f.close
    fok = open(filename,'w')
    fok.write(sub)
    fok.close()

def file_replace_escap(filename,srcStr,dstStr):
    f = open(filename,'r')
    sc = f.read()
    stra0 = re.sub('\^','\\^',srcStr)
    stra1 = re.sub('\$','\\$',stra0)
    stra2 = re.sub('\(','\\(',stra1)
    stra3 = re.sub('\)','\\)',stra2)
    stra = re.compile(stra3)
    sub = re.sub(stra,dstStr,sc)
    f.close
    fok = open(filename,'w')
    fok.write(sub)
    fok.close()

#www.iplaypy.com
def fetch_conf(fn,conf):  #conf maybe: IPADDR|NETMASK|NETWORK|GATEWAY
    #file_nw = file('/etc/sysconfig/network','r')
    #file_nw.seek(0,0)
    fn.seek(0,0)
    if conf == 'GATEWAY':
        #file_nw = file('/etc/sysconfig/network','r')
        for eachLine in fn:
            #alist = eachLine.strip().split('=')
            m1 = re.match('^#|^$',eachLine)
            if m1 is not None:
                continue
            else:
                m1 = re.search('=',eachLine)
                if m1 is not None:
                    alist = eachLine.strip().split('=')
                    nw[alist[0]]=alist[1]
        #print nw[conf]
        return  nw[conf]
    else:
        #file_eth1_pu = file('/etc/sysconfig/network-scripts/ifcfg-eth0','r')
        for eachLine in fn:
            #alist = eachLine.strip().split('=')
            m1 = re.match('^#|^$',eachLine)
            if m1 is not None:
                continue
            else:
                m1 = re.search('=',eachLine)
                if m1 is not None:
                    alist = eachLine.strip().split('=')
                    eth[alist[0]]=alist[1]
        return eth[conf]

def check_pn():
    file_eth1_pu = open(file_eth1,'r')
    xml_ip = parse_xml(filexml,"network","public_setting","ip")
    xml_mk = parse_xml(filexml,"network","public_setting","mask")
    xml_gw = parse_xml(filexml,"network","public_setting","gateway")
    conf_ip = fetch_conf(file_eth1_pu,"IPADDR")
    conf_mk = fetch_conf(file_eth1_pu,"NETMASK")
    conf_gw = fetch_conf(file_eth1_pu,"GATEWAY")
    file_eth1_pu.close()
    '''check public ip'''
    if  xml_ip[0] == conf_ip:
        print "ip is ok",xml_ip[0],conf_ip
    else:
        print "ip is wrong",xml_ip[0],conf_ip
        file_replace(file_eth1,conf_ip,xml_ip[0])
        #os.system('''sudo sed -i 's|wrong|right|g' /etc/sysconfig/network-scripts/ifcfg-eth1''')
    '''check public netmask'''
    if xml_mk[0] == conf_mk:
        print "netmask is ok",xml_mk[0],conf_mk
    else:
        print "netmask is wrong",xml_mk[0],conf_mk
        file_replace(file_eth1,conf_mk,xml_mk[0])
    '''check public gateway'''
    if xml_gw[0] == conf_gw:
        print "gateway is ok",xml_gw[0],conf_gw
    else:
        print "gateway is wrong",xml_gw[0],conf_gw
        file_replace(file_eth1,conf_gw,xml_gw[0])

def check_in():
    file_eth0_it = open(file_eth0,'r')
    xml_ip = parse_xml(filexml,"network","private_setting","ip")
    xml_mk = parse_xml(filexml,"network","private_setting","mask")
    xml_gw = parse_xml(filexml,"network","private_setting","gateway")
    conf_ip = fetch_conf(file_eth0_it,"IPADDR")
    conf_mk = fetch_conf(file_eth0_it,"NETMASK")
    conf_gw = fetch_conf(file_eth0_it,"GATEWAY")
    file_eth0_it.close()
    '''check internal ip'''
    if  xml_ip[0] == conf_ip:
        print "ip is ok",xml_ip[0],conf_ip
    else:
        print "ip is wrong",xml_ip[0],conf_ip
        file_replace(file_eth0,conf_ip,xml_ip[0])
    '''check internal netmask'''
    if xml_mk == conf_mk:
        print "netmask is ok",xml_mk[0],conf_mk
    else:
        print "netmask is wrong",xml_mk[0],conf_mk
        file_replace(file_eth0,conf_mk,xml_mk[0])
    '''check internal gateway'''
    if xml_gw == conf_gw:
        print "gateway is ok",xml_gw[0],conf_gw
    else:
        print "gateway is wrong",xml_gw[0],conf_gw
        file_replace(file_eth0,conf_gw,xml_gw[0])

def bu_cv(str):
    if str == "true":
       return str.replace('true','1')
    else:
       return str.replace('false','0')

def bu_cv0(str):
    if str == '0':
        return str.replace('0','true')
    else:
        return str.replace('1','false')
def check_sysctl_fs():
    file_sysctl_fs = open(file_sysctl,'r')
    xml_fs_en = "tcp_fast.enable=" + bu_cv(str(parse_xml(filexml,"fasttcp","enable")[0]))
    xml_fs_p1 = "tcp_fast.noaccel_lport=" + parse_xml(filexml,"fasttcp","block_port","value")[0]
    xml_fs_p2 = "tcp_fast.noaccel_lport2=" + parse_xml(filexml,"fasttcp","block_port","value")[1]
    conf_fs_en = "tcp_fast.enable=" + fetch_conf(file_sysctl_fs,'tcp_fast.enable')
    conf_fs_p1 = "tcp_fast.noaccel_lport=" + fetch_conf(file_sysctl_fs,'tcp_fast.noaccel_lport')
    conf_fs_p2 = "tcp_fast.noaccel_lport2=" + fetch_conf(file_sysctl_fs,'tcp_fast.noaccel_lport2')
    file_sysctl_fs.close()

    if xml_fs_en == conf_fs_en:
        print "fastsoft enable toggle is ok"
    else:
        print "fastsoft enable toggle is bad"
        file_replace(file_sysctl,conf_fs_en,xml_fs_en)

    if xml_fs_p1 == conf_fs_p1:
        print "fastsoft p1 is ok"
    else:
        print "fastsoft p2 is bad"
        file_replace(file_sysctl,conf_fs_p1,xml_fs_p1)

    if xml_fs_p2 == conf_fs_p2:
        print "fastsoft p2 is ok"
    else:
        print "fastsoft p2 is bad"
        file_replace(file_sysctl,conf_fs_p2,xml_fs_p2)

def check_cache():
    file_vd_cf = open(file_vd,'r')
    xml_cache_mm = "VARNISH_STORAGE_SIZE=" + parse_xml(filexml,"cache","max_memeory")[0] + "M"
    xml_cache_en = str(parse_xml(filexml,"cache","enable")[0])
    conf_cache_mm = "VARNISH_STORAGE_SIZE=" + fetch_conf(file_vd_cf,'VARNISH_STORAGE_SIZE')
    conf_cache_en = bu_cv0(str(os.system('ps aux | grep -v grep | grep -q varnishd')))
    file_vd_cf.close()
    if xml_cache_mm[0] == conf_cache_mm:
        print "conf_cache_mm is ok"
    else:
        print "conf_cache_mm is bad"
        file_replace(file_vd,conf_cache_mm,xml_cache_mm[0])
    if xml_cache_en[0] == conf_cache_en:
        print "conf_cache_en is ok"
    else:
        print "conf_cache_en is bad"
        if xml_cache_en == 'true':
            os.system('/sbin/service varnishd start')
        if xml_cache_en == 'false':
            os.system('/sbin/service varnishd stop')

def gen_backend_str():
    k = ""
    xml_varnish_origin = parse_xml(filexml,"varnish","origin","value")
    #print xml_varnish_origin[0].split('|')[0],xml_varnish_origin[0].split('|')[1],xml_varnish_origin[0].split('|')[2]
    ft_bk = "#origin_backend_start#" + "\n" + "backend  " + xml_varnish_origin[0].split('|')[0] + " {\n" + "\t.host = \"" + xml_varnish_origin[0].split('|')[1] + "\";\n\t.port = \"" + xml_varnish_origin[0].split('|')[2]  + "\";\n}\n"
    for i in xml_varnish_origin[1:-1]:
        #print i.split('|')[0],i.split('|')[1],i.split('|')[2]
        j = "backend  " + i.split('|')[0] + " {\n" + "\t.host = \"" + i.split('|')[1] + "\";\n\t.port = \"" + i.split('|')[2]  + "\";\n}\n"
        k += j
    #print xml_varnish_origin[-1].split('|')[0],xml_varnish_origin[-1].split('|')[1],xml_varnish_origin[-1].split('|')[2]
    lt_bk = "backend  " + xml_varnish_origin[-1].split('|')[0] + " {\n" + "\t.host = \"" + xml_varnish_origin[-1].split('|')[1] + "\";\n\t.port = \"" + xml_varnish_origin[-1].split('|')[2]  + "\";\n}\n#origin_backend_end#\n"
    all_str= ft_bk + k + lt_bk
    #print all_str
    return all_str


def gen_origin_str():
    k = ""
    xml_varnish_origin = parse_xml(filexml,"varnish","origin","value")

    ft_o = "#origin_start#" + "\n" + "if (req.http.host ~ \"^" + xml_varnish_origin[0].split('|')[0]  +  "$\") {\n" + "set req.backend=" + xml_varnish_origin[0].split('|')[0] + " ;\n}"
    #print ft_o
    #print xml_varnish_origin[1:]
    for i in xml_varnish_origin[1:]:
       #print i
       j = "  elsif (req.http.host ~ \"^" + i.split('|')[0]  +  "$\") {\n" + "set req.backend=" + i.split('|')[0] + " ;\n}"
       k += j
    lt_o = "  else {\n" + "error 404 \"Unknown virtual host\";\n" + "}\n#origin_end#\n"
    #print ft_o + k + lt_o
    return ft_o + k + lt_o



def check_varnishd():
    file_vd_cf = open(file_vd,'r')
    xml_varnish_status = parse_xml(filexml,"varnish","status")[0]
    xml_varnish_storage = parse_xml(filexml,"varnish","storage")[0]
    #xml_varnish_origin = parse_xml(filexml,"varnish","origin","value")
    conf_varnish_status = bu_cv0(str(os.system('ps aux | grep -v grep | grep -q varnishd')))
    conf_varnish_storage = "VARNISH_STORAGE_SIZE=" + fetch_conf(file_vd_cf,'VARNISH_STORAGE_SIZE')
    file_vd_cf.close()

    fl_vf = open(file_vf,'r')
    fl_vf_all = fl_vf.read()
    m1 = re.search('#origin_backend_start#(\n.*)*#origin_backend_end#\n',fl_vf_all)
    if m1 is not None:
        conf_varnishd_backend = m1.group()
    fl_vf.close()


    fl_vf2 = open(file_vf,'r')
    fl_vf_all2 = fl_vf2.read()
    m2 = re.search('#origin_start#(\n.*)*#origin_end#\n',fl_vf_all2)
    if m2 is not None:
        conf_varnishd_origin = m2.group()
    fl_vf2.close()


    xml_varnishd_backend = gen_backend_str()
    xml_varnishd_origin = gen_origin_str()

    if xml_varnishd_backend != conf_varnishd_backend:
        print "conf_varnishd_backend is wrong, Begin replace the wrong conf..."
        print xml_varnishd_backend
        print conf_varnishd_backend
        file_replace('/usr/local/varnish_214/etc/varnish_01.conf',conf_varnishd_backend,xml_varnishd_backend)
        #os.system('service varnishd restart')
    else:
        print xml_varnishd_backend
        print conf_varnishd_backend
        print "conf_varnishd_backend is ok"
    if xml_varnishd_origin != conf_varnishd_origin:
        print xml_varnishd_origin
        print '------------------'
        print conf_varnishd_origin
        print "conf_varnishd_origin is wrong, Begin replace the wrong conf..."
        file_replace_escap('/usr/local/varnish_214/etc/varnish_01.conf',conf_varnishd_origin,xml_varnishd_origin)
        #os.system('service varnishd restart')
    else:
        print xml_varnishd_origin
        print '------------------'
        print conf_varnishd_origin
        print "conf_varnishd_origin is ok"


def main():
    while True:
        time.sleep(3)
        daemonize()
        check_pn()
        check_in()
        check_sysctl_fs()

2000
   check_cache()
        check_varnishd()
        gen_backend_str()
        gen_origin_str()
        check_varnishd()

if __name__ == '__main__':
   main()
```


* 以上操作经过验证可以直接拿去。
* Owner `breakEval13`
* https://github.com/breakEval13