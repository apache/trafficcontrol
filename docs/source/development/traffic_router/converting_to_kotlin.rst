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

.. _tr-to-kotlin:

***********************************
Converting Traffic Router to Kotlin
***********************************


Kotlin Setup
============

#. Get the latest version of IntelliJ IDEA Community Edition from the package manager of your choice.

	.. Note:: Version 2022.2.2 or higher is recommended.

   Go to File -> Settings -> Plugins and update the Kotlin plugin if an update is available.

#. Remove the ``.idea/`` directory from the :file:`traffic_router/` router directory

#. Open the :file:`traffic_router/` directory in IntelliJ IDEA. It will auto-detect Traffic Router's Java version and Maven modules and download the Maven dependencies. In IntelliJ, the *Project* pane should look like this:

	.. figure:: intellij-traffic_router.png
		:width: 60%
		:align: center
		:alt: Traffic Router in IntelliJ Idea

		Project Pane

#. Go to Settings -> Editor -> Code Style -> Kotlin -> Imports and
	* Under *Top-Level Symbols*, choose *Use single name import*

	* Under *Java Statics and Enum Members*, choose *Use single name imports*

	* Under *Packages to Use Import with '*'*, remove packages until it says *Nothing to show*
	
	* Save your settings

Converting a single source file to Kotlin
=========================================

#. In the *Project* pane, right click on the main ``traffic_router`` directory and choose *Convert Java File to Kotlin File*. This process will take 20-30 minutes. Do not do anything else in IntelliJ until this step completes.

	.. Note:: It is very important to select the ``traffic_router`` directory so that *all* Java sources are considered during the conversion process. If you select only one or some Traffic Router sources before clicking *Convert Java File to Kotlin File*, the converter will incorrectly detect variables as nullable because not all sources were considered. Any Kotlin sources generated without selecting the whole ``traffic_router`` directory should be scrapped, even if you spent a lot of time on it.

#. Decide which single newly-converted Kotlin source file you want to keep. You may want to choose a small source file to start off with. To find the smallest newly-converted Kotlin source file:

	.. code-block:: shell
		:caption: List of newly-converted Kotlin source files, ordered by size

		git diff --cached --name-only '*.kt' | xargs wc -l | sort -gr

   In this example, we will keep the conversion of :file:`traffic_router/connector/src/main/java/org/apache/traffic_control/traffic_router/utils/HttpsProperties.java`.

#. Restore all of the other newly-deleted Java sources:

	.. code-block:: shell
		:caption: Restore all of the other newly-deleted Java sources

		git diff --cached --name-only '*.java' |
			grep -v 'traffic_router/connector/src/main/java/org/apache/traffic_control/traffic_router/utils/HttpsProperties\.*' |
			xargs git checkout HEAD

#. Delete all of the other newly-generated Kotlin sources:

	.. code-block:: shell
		:caption: Delete all of the other newly-generated Kotlin sources

		git diff --cached --name-only '*.kt' |
			grep -v 'traffic_router/connector/src/main/java/org/apache/traffic_control/traffic_router/utils/HttpsProperties\.*' |
			xargs git rm -f

#. Verify that running ``git status`` shows you only 2 changes (both of which are staged):

	* ``HttpsProperties.java`` is deleted
	
	* ``HttpsProperties.kt`` is added

	.. code-block:: shell
		:caption: The result of ``git status``

		Changes to be committed:
		  (use "git restore --staged <file>..." to unstage)
			deleted:    traffic_router/connector/src/main/java/org/apache/traffic_control/traffic_router/utils/HttpsProperties.java
			new file:   traffic_router/connector/src/main/java/org/apache/traffic_control/traffic_router/utils/HttpsProperties.kt

#. Commit these staged changes with a descriptive *commit* message.

#. Open your new Kotlin source file, most likely a class, in IntelliJ and look at it yourself before going on.

#. The conversion process mangles some imports. The most common example of this is having the word ``import`` at the end of one line, then the package name on a completely different line below:

	.. code-block:: kotlin
		:caption: An example of imports mangled from the Kotlin conversion process

		package org.apache.traffic_control.traffic_router.utilsimport

		import org.apache.logging.log4j.LogManager
		import org.apache.traffic_control.traffic_router.utils.HttpsProperties
		import java.nio.file.*
		import java.util.function.Consumer

		org.springframework.web.bind.annotation .RequestMapping

	Run this code snippet to fix these types of mangled imports:

	.. code-block:: shell
		:caption: Fix a common type of mangled import

		sed -i -z 's|import\(\n.*\n\)\([a-hk-z][a-z]*\.\)|\1import \2|g' $(git ls-files '*.kt')


	Although the snippet above fixes that specific type of mangled import, in cases where the conversion mistakenly puts multiple imports on a single line, you will need to fix those yourself.

#. Once the imports are syntactically correct, remove the unused imports. With your Kotlin source file open in IntelliJ, go to Code -> Optimize Imports. This will remove a lot, likely hundreds, of unused imports left over from the Kotlin conversion process.

   Stage and commit these fixes to the imports before going on.

#. Familiarize yourself with Kotlin syntax, particularly with syntax around null safety (go through https://kotlinlang.org/docs/null-safety.html).

   A big Java pain point is null checks and NullPointerExceptions, and one of Kotlin's greatest features is its compile-time null safety. Converting your class from Java to Kotlin most likely produced several compile errors, which are usually null safety-related.

   Jump to the first compile error by pressing ``F2``. In the example of ``HttpsProperties.kt``:

	.. figure:: null-compile-error.png
		:width: 80%
		:align: center
		:alt: The first compile-time error in HttpsProperties.kt

		Compile-time error in HttpsProperties.kt


   In this case, the null safety errors and warnings can be eliminated by making ``HTTPS_PROPERTIES_FILE`` a ``const val`` (in Java, ``HTTPS_PROPERTIES_FILE`` was ``private static final``):

	.. figure:: errors-warnings-fixed.png
		:width: 80%
		:align: center
		:alt: HttpsProperties.kt with errors and warnings fixed

		HttpsProperties.kt with errors and warnings fixed

   Almost all of these errors and warnings can be fixed with only small changes like these. Fix all of the Kotlin errors and warnings in the source file, then stage and commit.

#. Run unit, integration, and feature tests. Resolve test failures without removing the failing test cases.

#. PR your class! Congratulations, you have made the Traffic Router codebase safer and more maintainable.

	.. Note:: For very small source files, sometimes it is acceptable to include more than one converted source file in a single PR. However, any large converted source file should have its own PR.

To see the git history of the HttpsProperties example shown here, see https://github.com/zrhoffman/trafficcontrol/commits/tr-kotlin-HttpsProperties.
