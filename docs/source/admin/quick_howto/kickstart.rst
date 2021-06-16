..
..
.. Licensed under the Apache License, Version 2.0 (the "License");
.. you may not use this file except in compliance with the License.
.. You may obtain a copy of the License at
..
..     http://www.apache.org/licenses/LICENSE-2.0
..
.. Unless required by applicable law or agreed to in writing, software
.. distributed under the License is distributed on an "AS IS" BASIS,
.. WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
.. See the License for the specific language governing permissions and
.. limitations under the License.
..

.. _Creating-CentOS-Kickstart:

**********************************
Creating the CentOS Kickstart File
**********************************
The Kickstart file is a text file, containing a list of items, each identified by a keyword. This file can be generated using the `Red Hat Kickstart Configurator application <https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/5/html/installation_guide/ch-redhat-config-kickstart>`_, or it can be written from scratch. The Red Hat Enterprise Linux installation program also creates a sample Kickstart file based on the options selected during installation. It is written to the file :file:`/root/anaconda-ks.cfg` in this case. This file is editable using most text editors.

Generating a System Image
=========================
#. Create a Kickstart file.
#. Create a boot media with the Kickstart file or make the Kickstart file available on the network.
#. Make the installation tree available.
#. Start the Kickstart installation.

.. code-block:: shell
	:caption: Creating a New System Image Definition Tree from an Existing One

	# Starting from the Kickstart root directory (`/var/www/files` by default)
	mkdir newdir
	cd newdir/

	# In this example, the pre-existing system image definition tree is for CentOS 7.4 located in `centos74`
	cp -r ../centos74/* .
	vim ks.src
	vim isolinux/isolinux.cfg
	cd ..
	vim osversions.json

:file:`ks.src` is a standard, Kickstart-formatted file that the will be used to create the Kickstart (ks.cfg) file for the install whenever a system image is generated from the source tree. :file:`ks.src` is a template - it will be overwritten by any information set in the form submitted from :menuselection:`Tools --> Generate ISO` in Traffic Portal. Ultimately, the two are combined to create the final Kickstart file (:file:`ks.cfg`).

.. note:: It is highly recommended for ease of use that the system image source trees be kept under 1GB in size.

.. seealso:: For in-depth instructions, please see `Kickstart Installation <https://access.redhat.com/documentation/en-US/Red_Hat_Enterprise_Linux/6/html/Installation_Guide/s1-kickstart2-howuse.html>`_ in the Red Hat documentation.

.. _kickstart.files.location:

``kickstart.files.location``
=============================

The Kickstart root directory used by :ref:`to-overview` (``/var/www/files`` by default) can be changed by setting the ``kickstart.files.location`` :term:`Parameter`.
