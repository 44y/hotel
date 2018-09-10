# -*- coding: UTF-8 -*-

import json, os, shutil, stat, subprocess

#生成exe版本号
VERSION_Major = 1
VERSION_Minor = 0
VERSION_Patch = 0
VERSION_Build = 2

def readJson(filename):
    with open('versioninfo.json', 'r') as load_f:
        load_dict = json.load(load_f)
        return load_dict

def writeJson(new_dict,filename):
    with open(filename, 'w') as dump_f:
        json.dump(new_dict, dump_f)


JSON_NAME='versioninfo.json'

got_dict=readJson(JSON_NAME)

Major=str(VERSION_Major)
Minor=str(VERSION_Minor)
Patch=str(VERSION_Patch)
Build=str(VERSION_Build)

got_dict["FixedFileInfo"]["FileVersion"]["Major"] = VERSION_Major
got_dict["FixedFileInfo"]["FileVersion"]["Minor"] = VERSION_Minor
got_dict["FixedFileInfo"]["FileVersion"]["Patch"] = VERSION_Patch
got_dict["FixedFileInfo"]["FileVersion"]["Build"] = VERSION_Build

tmp_id=os.popen('git log -1 --pretty=%h').readlines()
commit_id = tmp_id[0].strip('\n')
#print(commit_id)

new_version = Major+"."+Minor+"."+Patch+"."+Build+"."+commit_id
got_dict['StringFileInfo']['ProductVersion'] = new_version
writeJson(got_dict, JSON_NAME)



output_path="ik_audit_server"
bin_path = output_path+"/bin/"
isExist = os.path.exists(bin_path)
if isExist == False:
    os.makedirs(bin_path, stat.S_IWOTH)


exe_file="IK_Auth_Srv.exe"
os.system("go clean")
os.system("go generate")
os.system("go build -o "+exe_file)


#if os.path.exists(bin_path + exe_file):
os.remove(bin_path + exe_file)
shutil.move(exe_file, bin_path)

#if os.path.exists(bin_path + "conf.json"):
os.remove(bin_path + "conf.json")
shutil.copy("conf.json", bin_path)

#if os.path.exists(output_path +"/install.bat"):
os.remove(output_path +"/install.bat")
shutil.copy("../install/install.bat", output_path +"/")

#if os.path.exists(output_path +"/Readme.txt"):
os.remove(output_path +"/Readme.txt")
shutil.copy("../install/Readme.txt", output_path +"/")

'''
if os.path.exists(output_path +"/uninst.bat"):
    os.remove(output_path +"/uninst.bat")
shutil.copy("../install/uninst.bat", output_path +"/")
'''

#使用powershell的Compress-Archive进行压缩
#if os.path.exists(output_path+".zip"):
os.remove(output_path+".zip")
zip_cmd_argv = ["powershell", r"Compress-Archive", r"-Path", output_path, r"-DestinationPath", output_path+".zip"]
subprocess.call(zip_cmd_argv)
