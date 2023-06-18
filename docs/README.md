<!--
    Licensed to the Apache Software Foundation (ASF) under one
    or more contributor license agreements.  See the NOTICE file
    distributed with this work for additional information
    regarding copyright ownership.  The ASF licenses this file
    to you under the Apache License, Version 2.0 (the
    "License"); you may not use this file except in compliance
    with the License.  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing,
    software distributed under the License is distributed on an
    "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
    KIND, either express or implied.  See the License for the
    specific language governing permissions and limitations
    under the License.
-->

# SPHINX Documentation Local Build Guide 
<details>
<summary>Windows</summary>

Requirements: 
* [Python](https://www.python.org/downloads/windows/)
* [Make](http://gnuwin32.sourceforge.net/packages/make.htm)
* Powershell 'Run as Admininstrator'

Steps:
* Verify Python Installation
  ```
  PS C:\Users\Administrator> python --version
  Python 3.11.4
  ```
* Verify Make installation
  ```
  PS C:\Users\Administrator> make --version
  GNU Make 3.81
  Copyright (C) 2006  Free Software Foundation, Inc.
  This is free software; see the source for copying conditions.
  There is NO warranty; not even for MERCHANTABILITY or FITNESS FOR A
  PARTICULAR PURPOSE.
  
  This program built for i386-pc-mingw32
  ```
* Create a new python virtual environment
  ```
  PS C:\Users\Administrator\Downloads> python -m venv docs_env
  ```
* Activate the Virtual Environment
  ```
  PS C:\Users\Administrator\Downloads> .\docs_env\Scripts\Activate.ps1
  ```
* Navigate to docs directory in trafficcontrol repo and install the requirements
  ```
  (docs_env)PS C:\Users\Administrator\Downloads\trafficcontrol\docs> pip install -r source/requirements.txt
  ```
* Upgrade Sphinx and Sphinx-autobuild
  ```
  (docs_env)PS C:\Users\Administrator\Downloads\trafficcontrol\docs> pip install --upgrade sphinx sphinx-autobuild
  ```
* Build docs files locally with necessary make commands
  ```
  (docs_env)PS C:\Users\Administrator\Downloads\trafficcontrol\docs> make html
  ```
  ```
  (docs_env)PS C:\Users\Administrator\Downloads\trafficcontrol\docs> make clean
  ```

</details>

<details>

<summary>Mac/Linux</summary>

Requirements:
* Python
* Make (Usually comes pre-installed)

Steps:
* Create a new python virtual environment
  ```
  python -m venv new_rst
  ```
* Activate the Virtual Environment
  ```
  source path_to_your_virtual_env/new_rst/bin/activate
  ```
* Navigate to docs directory in trafficcontrol repo and install the requirements
  ```
  pip install -r requirements.txt
  ```
* Build docs files locally with necessary make commands
  ```
  make html
  ```
  ```
  make clean
  ```
</details>
