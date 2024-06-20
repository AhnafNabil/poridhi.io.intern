# Setting Up a Docker Swarm Cluster: Installing a Hypervisor**


**Introduction:**
In this tutorial, we will guide you through setting up a Docker Swarm cluster by installing a Hypervisor on your host machine. A Hypervisor is a piece of software that allows us to create virtual machines, which will be used to host our Docker Swarm cluster.

**Step 1: Updating sources.list file**
Open the sources.list file, which contains repository links for software updates. Add the necessary repository link for VirtualBox.

```bash

sudo nano /etc/apt/sources.list
```

Add the following line to the file:

```
deb [arch=amd64] https://download.virtualbox.org/virtualbox/debian bionic contrib
```

Save the file and exit the text editor.

**Step 2: Getting the GPG key for VirtualBox**

Retrieve the GPG key for VirtualBox to ensure the authenticity of the downloaded packages.

```bash
wget -q https://www.virtualbox.org/download/oracle_vbox_2016.asc -O- | sudo apt-key add -
```

**Step 3: Updating apt package manager**
Update the apt package manager to include the newly added VirtualBox repository.

```bash
sudo apt-get update
```

**Step 4: Installing VirtualBox**
Install VirtualBox using the apt package manager. Specify the version of VirtualBox to install (e.g., VirtualBox 5.2).

```bash
sudo apt-get install virtualbox-5.2
```

Follow the on-screen instructions to complete the installation process.

**Step 5: Verifying VirtualBox Installation**

Check if VirtualBox has been successfully installed by searching for it in the list of installed software.

```bash
# Open the application menu or use the search functionality to find VirtualBox
```
Verify that Oracle VirtualBox appears in the list of installed software, indicating a successful installation.

![virtualbox](./images/4.png)

**Conclusion:**
Congratulations! We have successfully installed VirtualBox, a Hypervisor, on your host machine. You are now ready to proceed with setting up a Docker Swarm cluster using virtual machines.

