/*
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package configuration;

import org.apache.traffic_control.traffic_router.neustar.configuration.NeustarConfiguration;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.test.context.ContextConfiguration;
import org.springframework.test.context.junit4.SpringJUnit4ClassRunner;
import org.springframework.test.context.support.AnnotationConfigContextLoader;

import java.io.File;
import java.nio.file.Path;
import java.nio.file.Paths;

import static org.hamcrest.CoreMatchers.equalTo;
import static org.hamcrest.MatcherAssert.assertThat;

@RunWith(SpringJUnit4ClassRunner.class)
@ContextConfiguration(classes={NeustarConfiguration.class, PeriodicUpdateConfigurationTest.TestConfiguration.class}, loader=AnnotationConfigContextLoader.class)
public class PeriodicUpdateConfigurationTest {
	@Autowired
	File neustarDatabaseDirectory;

	@Autowired
	File neustarOldDatabaseDirectory;

	@Test
	public void itUsesCorrectDirectoriesForDataDownload() {
		assertThat(neustarDatabaseDirectory.getAbsolutePath(), equalTo("/opt/traffic_router/db/neustar"));
		assertThat(neustarOldDatabaseDirectory.getAbsolutePath(), equalTo("/opt/traffic_router/db/neustar/old"));
	}

	@Configuration
	static class TestConfiguration {
		@Bean
		public Path databasesDir() {
			return Paths.get("/opt/traffic_router/db");
		}
	}
}