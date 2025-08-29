# Release Package Manager GUI Application
## By Arnav Josh
This program was written as a part of my internship with **Globus Medical**.
Globus Medical is a medical device company that develops and manufactures products for musculoskeletal disorders, including spinal and orthopedic implants, and various . I worked with their software deployment team - which works on making software updates to their robots quickly and effectively. 
One deployment tool is the **Release Package Manager**, which creates, signs, and publishes to the cloud a **manifest file**.
The manifest file is a json that holds information like the software versions and builds, the robot which is being updated, and install methods/orders for **artifacts**. Artifacts are files that go along with the software update.
The RPM tool was only in a command line interface previously, so I was enlisted to create a GUI application for the company.
The code was written in Go with the Fyne gui library

**For the purposes of submitting this as my summer project, I have removed references to the companies private internal code, making it into a demo gui app**
## Installation
By running build.ps1, an exe file with the program is created. However, this requires the GCC compiler and fyne library installs, so I have provided the exe file itself. 
**To run the program, run the file "rmpg.exe" in the exe folder**
## How to run the program
If you want to create a new manifest, then you must enter inputs into the fields provided (ex. platform, version, build)
It is also important to select the artifact folder (which contains all the artifacts) and arrange/edit them according to your needs
Then, you can sign a manifest (either one that has just been created, or you can load one in), with an API server selected
The next steps are to stage and publish the manifest, which will put them onto the cloud

However, these will not work in the demo version I've provided, because it uses private company code.
## Interpreting output
The manifest will be saved locally at a location that the user can select
It will also be sent to the company's cloud system depending on whether or not the user completes the stage and publish steps

However, in this demo version, there will not be any output in the files of the user, only dialog messages sent to the user.
## UML Diagram
<img width="766" height="442" alt="{91BF4F62-3C10-481F-AEE7-04F3C4B39772}" src="https://github.com/user-attachments/assets/92a730e0-65b0-46ad-a9a6-ecbdf0731fdc" />

