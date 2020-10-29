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

# IDE Setup

Most all of the code contained uses Builders for implementation. This ensures all properties are managed and defaults are used. As well as facilitate YAML based construction. 
Most all of these builders are constructed using Google's AutoValue library. This allows for auto code generation based on annotations and abstract classes. 
To support this within your IDE you will need to do a couple things listed below.

## Eclipse

1. Install m2e-apt from the Eclipse marketplace. Help -> Eclipse Marketplace -> Search "m2 apt" -> Install m2e-apt
2. Activate the apt processing. Preferences -> Maven -> Annotation processing -> Switch to Experimental
3. Import project or if it has already been imported refresh the projects form the maven sub-menu.

## IntelliJ

1. Open Annotation Processors settings. Settings -> Build, Execution, Deployment -> Compiler -> Annotation Processors
2. Select the following buttons:
   * Enable annotation processing
   * Obtain processors from project classpath
   * Store generated sources relative to: Module content root
3. Set the generated source directories to be equal to the Maven directories:
   * Set “Production sources directory:” to t"arget/generated-sources/annotations"
   * Set “Test sources directory:” to "target/generated-test-sources/test-annotations"