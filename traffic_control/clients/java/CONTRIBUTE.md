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