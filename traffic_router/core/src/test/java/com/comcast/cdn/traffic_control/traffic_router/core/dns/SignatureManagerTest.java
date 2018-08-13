/*
t *
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

package com.comcast.cdn.traffic_control.traffic_router.core.dns;

import com.comcast.cdn.traffic_control.traffic_router.core.edge.CacheRegister;
import com.comcast.cdn.traffic_control.traffic_router.core.edge.InetRecord;
import com.comcast.cdn.traffic_control.traffic_router.core.config.ConfigHandler;
import com.comcast.cdn.traffic_control.traffic_router.core.config.SnapshotEventsProcessor;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryService;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.SteeringWatcher;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.AnonymousIpDatabaseService;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.FederationRegistry;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.FederationsWatcher;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.MaxmindGeolocationService;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker;
import com.comcast.cdn.traffic_control.traffic_router.core.router.TrafficRouter;
import com.comcast.cdn.traffic_control.traffic_router.core.router.TrafficRouterManager;
import com.comcast.cdn.traffic_control.traffic_router.core.util.JsonUtils;
import com.comcast.cdn.traffic_control.traffic_router.core.util.TrafficOpsUtils;
import com.comcast.cdn.traffic_control.traffic_router.geolocation.GeolocationService;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.google.common.cache.LoadingCache;
import org.apache.commons.io.IOUtils;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.api.mockito.PowerMockito;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;
import org.powermock.reflect.Whitebox;
import org.xbill.DNS.RRset;
import org.xbill.DNS.Record;
import org.xbill.DNS.Type;
import org.xbill.DNS.Zone;

import java.io.File;
import java.io.IOException;
import java.io.InputStream;
import java.util.Iterator;
import java.util.List;
import java.util.Map;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.*;
import static org.junit.Assert.*;
import static org.mockito.Matchers.*;
import static org.mockito.Mockito.*;
import static org.powermock.api.mockito.PowerMockito.*;

class SigManagerForTesting extends SignatureManager {
	public static enum KeyProfile {
		ONE, TWO, THREE;
	}

	public static KeyProfile returnKey = KeyProfile.ONE;

	public SigManagerForTesting(final ZoneManager zoneManager, final CacheRegister cacheRegister,
	                            final TrafficOpsUtils trafficOpsUtils,
	                            final TrafficRouterManager trafficRouterManager) {
		super(zoneManager, cacheRegister, trafficOpsUtils, trafficRouterManager);
	};

	@Override
	protected JsonNode fetchKeyPairData(CacheRegister cacheRegister){
		final String cdnZskStr =
			"[" +
				"{" +
					"\"inceptionDate\": 1545136505," +
					"\"expirationDate\": 2547728505, " +
					"\"name\": \"thecdn.example.com.\"," +
					"\"ttl\": \"60\", " +
					"\"status\": \"new\"," +
					"\"effectiveDate\": 1543408205," +
					"\"public\": " +
					"\"OyBUaGlzIGlzIGEgem9uZS1zaWduaW5nIGtleSwga2V5aWQgNjUwMTYsIGZvciB0aGVjZG4uZXhhbXBsZS5jb20uCjsgQ3JlYXRlZDogMjAxOTA1MTkwMjQ2MzAgKFNhdCBNYXkgMTggMjA6NDY6MzAgMjAxOSkKOyBQdWJsaXNoOiAyMDE5MDUxOTAyNDYzMCAoU2F0IE1heSAxOCAyMDo0NjozMCAyMDE5KQo7IEFjdGl2YXRlOiAyMDE5MDUxOTAyNDYzMCAoU2F0IE1heSAxOCAyMDo0NjozMCAyMDE5KQp0aGVjZG4uZXhhbXBsZS5jb20uIElOIEROU0tFWSAyNTYgMyA4IEF3RUFBZUFYOU1NNDl6cm1uT3Vwc0haNzJTNlEyekdSRUtsTnM1M3AzMzFEeU1PUUx5TTV3OWRqIEgyS1pCZk9MSlk2V2dVY2VlWWdqSjNtL0JmMDZvZ0JKQVkvRXVMR1FzZVJmdXROSUhDZnhhOERYIE9PR1luaVlGa3pxSHhJcVRIVVhUemVFNjh5ZVAvS21pckRwOWRMRmxQdUtIZFN5RVVxbGthdUFNIExUYWY2V1k5TVlsWkYwelV1amR4NTJUbXhQcnR2YU9DQjZTWTFsZnNCdFdsc2tsZ210ZTk3TnE3IDNYaW8rS0MxejNJT0lGMHBzOG5TRU9udDlNdTJEV2lEa1NXYkF4cnVpZkNDcWhTeThuNCtobUwzIHl2NEdTWnpqYUZQNHlYTjRxM3RKNWVPTllBTnVULzlJK0VHY0pGZFVtOUVOV1ZPV2lzRlRNWDljIHpJbndWakxMNnlrPQo=\"," +
					"\"private\": " +
					"\"UHJpdmF0ZS1rZXktZm9ybWF0OiB2MS4zCkFsZ29yaXRobTogOCAoUlNBU0hBMjU2KQpNb2R1bHVzOiA0QmYwd3pqM091YWM2Nm13ZG52WkxwRGJNWkVRcVUyem5lbmZmVVBJdzVBdkl6bkQxMk1mWXBrRjg0c2xqcGFCUng1NWlDTW5lYjhGL1RxaUFFa0JqOFM0c1pDeDVGKzYwMGdjSi9GcndOYzQ0WmllSmdXVE9vZkVpcE1kUmRQTjRUcnpKNC84cWFLc09uMTBzV1UrNG9kMUxJUlNxV1JxNEF3dE5wL3BaajB4aVZrWFROUzZOM0huWk9iRSt1MjlvNElIcEpqV1Yrd0cxYVd5U1dDYTE3M3MycnZkZUtqNG9MWFBjZzRnWFNtenlkSVE2ZTMweTdZTmFJT1JKWnNER3U2SjhJS3FGTEx5Zmo2R1l2ZksvZ1pKbk9Ob1UvakpjM2lyZTBubDQ0MWdBMjVQLzBqNFFad2tWMVNiMFExWlU1YUt3Vk14ZjF6TWlmQldNc3ZyS1E9PQpQdWJsaWNFeHBvbmVudDogQVFBQgpQcml2YXRlRXhwb25lbnQ6IGhTWTJ6SGhnamFMUXdrWEZEK0Z1bmZoOEJPdUkxcy9RVlpmTXJ3VkRtTTltOHFzejdocDdYMzVFeHZ4NmlFcHM3ZkR4alM0MkdtU2lEbUIwT1c4bmVtRU16SlNJM29UeTRPOExxOEFLb2o0L0VldzRxNmJNWUE2amZTTUVWYVVQd3kvZm9qWXZqaXVWNGZzQkJ4WjlVdVBoZjEzd0w2MzJ3c0Q1YkdkL3FycThvdk1nUWsrT1pOL2VMWVlyR1c4R2RBeVJQWlAxVTVZVmVQd2Jrc2s4Qy9NbkZTMk13aU84NFRlVmxXOU9obkZsZmRUNHhsWW8wMTljdEx0dWs2QTVGK2doOXB0RmlEaUN6RnVydE1jbGtZQy9EVGZWMzRzOExhcEtmeHBjK21aT0xBdVRqem1iNjhKL244aERXd05uTHhIZHVndzcyS3hjQnBGTGp5VTBjUT09ClByaW1lMTogOWorWXAzNzEreXIxbTBmc0M2VW5lMWJQYWVXUjNpYTh0SERHWDZPaWFCN1ROS0xEb3c0VHo0d0dqdzU3ZW1MTWVmL3puRS9adlIwMHc0dTVDcGxTVmdKWnNxZUFKTjFmbmRqZGd1UEpTeUgwaUZwRERRbTBRSVN1M080U3pNREpaMWdySVRLQUpCTWVIb0lIY0MrMS9lVlhNdHQxdDNQdVcwMUFGRzZiSExzPQpQcmltZTI6IDZQZkNkMy90YWxGRFphb2pnWjFYTG9mTEVGVjJ0a3VRZ2hPTlUxYWkvUkdsMUFNaWxHZ0FsdkZkWmJLUUVrOUFQenl3OGR1L3oyRFNwb2d4ZE1CT2xwWEYwZG1ic3Q3dHBha3p1dU5oVjVoSmdDbWdHWC9hZTBDaVhFaFdqSGdHVUlta3Z5MTRheXRTSEc3N0ovclB0QlpoK0prLzNmNDZvZ1padWVueXEycz0KRXhwb25lbnQxOiBkUUZhNDV0ci9lQTN1NFM5SC90bGV6R1FkRnprcG8zNWRETnh6dGZOdjNPR0g2aUhGTjhIZ3NLaUN5OXlkSUNsY0FLeUdqL0cwaGtpalJmQzRNcGhXUVRjNGdxODFNZHJPM0ZrTDJGVXNDSitNcHZRNkUwSUhPL3V0b05ZNDNsbk9YZW5acXE4UUdmbEU5SHMvdDFzdUN0VTk1RlhxY2VvSmpIYWxOckpVU2s9CkV4cG9uZW50MjogZjNvYk1qS1JJZzBIZVJqcFJ1Sm1zekpnL2pZVnVGaU80VzU3ZGEvQmpnTGJIN0Q3ZWdPUzF3V0oycjBqc3JFazhiVnZDVmo2a3lwaStmY3FRTDErQTA0L0tiNE9RWWxVdHNKb2FRcEUySjZmRVg4MWVONktHY2xiVG0zUzFoaFRORHl0Sm1ObU1oWkpYdituZE0vOFdZbXA0Rk5UTEZFUm1sa3BQZDd6QjdNPQpDb2VmZmljaWVudDogeFNGZjdKL1hPUWpMWSttRkNYSTFPc0h4Z0dlTGhBTitnbkJHU0JYNUVBT1pWT2tJNG5HcWlQNU1ZRWlONDRJc29oQ1lwejlhcDluUEhKUWJPVmFPaUFhdVJxT2V4RTFzaEdwNzVVR0w4cHJ0b2NkTFJhaS83Q1hKMUFlWmJMdGllR280ZXQvbzdzUHFNMUVlYndlTnU2c3ppZ0dzSzRuQmlNTEF5aDh2UjlBPQpDcmVhdGVkOiAyMDE5MDUxOTAyNDYzMApQdWJsaXNoOiAyMDE5MDUxOTAyNDYzMApBY3RpdmF0ZTogMjAxOTA1MTkwMjQ2MzAK\"" +
				"}" +
			"]";
		final String httpsZskStr =
			"[" +
				"{" +
					"\"inceptionDate\": 1545136505," +
					"\"expirationDate\": 2547728505, " +
					"\"name\": \"https-only-test.thecdn.example.com.\"," +
					"\"ttl\": \"60\", " +
					"\"status\": \"new\"," +
					"\"effectiveDate\": 1543408205," +
					"\"public\": " +
					"\"OyBUaGlzIGlzIGEgem9uZS1zaWduaW5nIGtleSwga2V5aWQgMjEyMzEsIGZvciBodHRwcy1vbmx5LXRlc3QudGhlY2RuLmV4YW1wbGUuY29tLgo7IENyZWF0ZWQ6IDIwMTkwNTE4MjIyODE2IChTYXQgTWF5IDE4IDE2OjI4OjE2IDIwMTkpCjsgUHVibGlzaDogMjAxOTA1MTgyMjI4MTYgKFNhdCBNYXkgMTggMTY6Mjg6MTYgMjAxOSkKOyBBY3RpdmF0ZTogMjAxOTA1MTgyMjI4MTYgKFNhdCBNYXkgMTggMTY6Mjg6MTYgMjAxOSkKaHR0cHMtb25seS10ZXN0LnRoZWNkbi5leGFtcGxlLmNvbS4gSU4gRE5TS0VZIDI1NiAzIDggQXdFQUFjNng0Mkd2aG4xUnk5Q1lzQnFoc3RaaURGMytrVktIaXpaZDlLYnJ1c2dRVlNlYVp3YWIgSkcvY3hucHU2WmpocnZmTFFzeFNzTDV0ZmpQRmV4cWEwcCtqTHJPVEJqWWVMYUx6VkpRcFJBYnYgTUxzN1VqQ3lEckZnRUt3WTQrNGhhV2NDcVl2Y01UN0xqYnRvTmxwYlNOb3o4YTFCTFFPMjVYdXogd29iNzlObWRFWVdhWEJsaEExVXVxWTVYUWhKbVBqTTE3a1lsU3cvTnVCczIzYkN6VUorVHVPcU4gTHljS2FWYlh3T0ZUZFNEM0FiUU5LbGxvUzY4NGdpRUxtUThjRHUwNlJwc3lqRk5XZEx2SC9wckwgMmtNbUNPVHJmUlJxdU0zNVZDWGkzQU5nekRqSzg2T0xuYzhiOEZoUDFiZWs2ZjVnWEc1NUlaZmggVXViVkpoSFdmUDg9Cg==\"," +
					"\"private\": " +
					"\"UHJpdmF0ZS1rZXktZm9ybWF0OiB2MS4zCkFsZ29yaXRobTogOCAoUlNBU0hBMjU2KQpNb2R1bHVzOiB6ckhqWWErR2ZWSEwwSml3R3FHeTFtSU1YZjZSVW9lTE5sMzBwdXU2eUJCVko1cG5CcHNrYjl6R2VtN3BtT0d1OTh0Q3pGS3d2bTErTThWN0dwclNuNk11czVNR05oNHRvdk5VbENsRUJ1OHd1enRTTUxJT3NXQVFyQmpqN2lGcFp3S3BpOXd4UHN1TnUyZzJXbHRJMmpQeHJVRXRBN2JsZTdQQ2h2djAyWjBSaFpwY0dXRURWUzZwamxkQ0VtWStNelh1UmlWTEQ4MjRHemJkc0xOUW41TzQ2bzB2SndwcFZ0ZkE0Vk4xSVBjQnRBMHFXV2hMcnppQ0lRdVpEeHdPN1RwR216S01VMVowdThmK21zdmFReVlJNU90OUZHcTR6ZmxVSmVMY0EyRE1PTXJ6bzR1ZHp4dndXRS9WdDZUcC9tQmNibmtobCtGUzV0VW1FZFo4L3c9PQpQdWJsaWNFeHBvbmVudDogQVFBQgpQcml2YXRlRXhwb25lbnQ6IGdCdGJEZEdIYnFiQ3h4L0xqa1lJSEwyYVlxdUVFSDIzOTlOVjdoc09vaThWU0QxM2UyMnVzTEpLUmRuMmFHNEFUczZwTVJCVFFVT0Z3N3F6S1BNcWRneml4aVBxNXFIZnNTNVZqcHJnOGRkUUFjdXJqa2pkdUd3TkNVMUgvL0N2Ymt3Rkg0MHk3SE9tL2ErQ2VKQTVZQnh2dnUwMVpDYjRwcW5zZUZNekY2Z1g2UDlGbnVMYXZKQ1JkLzBFTXJMZmRJbk5xMjkwTWhHT01nR09qRGxpWEcwVWdSWU5FZm9UWDlNL0R2dk9kek4vMDhDSEZhY0JtZUVSUW5WdWZmZUVUQXFoc0RUNVBLZTVUYnM0Y0JDbElQNFpEa3lWNWRrdCtHL1RPenZaZ0J1MzV5WlZhSW1jZDRlY3JMYWFvY0liSzUvUXhjS2NZakhPZFhaYk1uMjl1UT09ClByaW1lMTogNk15VlNQelFBN0k3R3pCMkc4U21RSEkraEhjUXlLK1QyY1ZDNFIwSWdBcHhheWtyWjVlMktiVGJZTmZva2JnTkRoUGlvQVZScU1Ra3k1K3J2RGxWeFRSTGh4ejk5MWJDUldmbWhOdFB2UGpLQUt6c3ppc3lFOG1DYm1DUS96WnVIVVF6Z0c3bm90ZjFPYjB2Q3pYMGd1WStLaXkrRzZjeHJ0RGhxRTc4VVJVPQpQcmltZTI6IDQwdE9KM2xQZjloVE9KbGxucnVydnlsTytuL2VzZC9TTkMyYk1hYTdqb2NPekV1YXJWOTVPUktRSEsvcEpKSXRKZjVQYXNpUjhVOGhlNkYwdXYwQVJueFhVNXU4RHB1VDZnUTdzQ1NPbDQrSEZYN2JGRkNrOE9vSWZqT3h4Tk9GSTkrVmY4RXhHVkUyaC8yb3ZVdFowRVBwSG1BSlpsbXI0d29YVDNIOFVzTT0KRXhwb25lbnQxOiBhQ1YyTjhQYUwzMGgwaVVaQVkwMUx5bXM3RWZ6KzBRSktlaUU1ZjYrR2NJU1FYV1dsRzZic1FiWENma3RjMXRhZzh6RG13RW1LaEV0d09hNnhxY0R4d2lCTFgzNzVCWXRMUzJ4UkFoMUlMNVFhSUwwSWJ2VFdHVFM5QnhCWWR0dzRhanNQVzNnTk1yV1N6Rm1oV1pxNzlDZlNQRGhUNis1bTFLWlVWbWNxWTA9CkV4cG9uZW50MjogazgydmQ5bmFDWitwbGwraUJsT3h2bkJsVEY0RWVaUzdnM1M0dTlQWm1UaFlOaDlmNlNmeGsxeHYzRWZFQ3lVSE9QS2p3Q1BIUzYwU2IrdXhGYnRhQjN2cDZaT1crY1RQcmRpczI4RVovSkszM0JHTzh5bng2RHUzNUNGSGsxK2M3NVFBQ29DZHBnSDZ2UG9GVlhyL1g0QVp3c2ZldFBEUTVxWDBQSlE1NmJNPQpDb2VmZmljaWVudDogNTBwNUpJVUwvdnI3eUhDeDl3a0JmUjRSMkVlTHMvZlFFdkRwR3c2ekg1Qng1Q0ZaLzhDcFQ0ejdVMWk5WjFkNHN1L0Y3WU5PWWdQMVI3SkZqYWxQcEY1OENZbHVZK1ZjUGEzT3NCdlZUUU94TE9jdWcza3lkZGtRVk03aHdScVhxamYwV0I3NENHejZoeUlrNFlhQWlYc0NtV0E4dFB1bFlwcUpBSURJYUtFPQpDcmVhdGVkOiAyMDE5MDUxODIyMjgxNgpQdWJsaXNoOiAyMDE5MDUxODIyMjgxNgpBY3RpdmF0ZTogMjAxOTA1MTgyMjI4MTYK\"" +
				"}" +
			"]";
		final String dnstZskStr =
			"[" +
				"{" +
					"\"inceptionDate\": 1545136505," +
					"\"expirationDate\": 2547728505, " +
					"\"name\": \"dns-test.thecdn.example.com.\"," +
					"\"ttl\": \"60\", " +
					"\"status\": \"new\"," +
					"\"effectiveDate\": 1543408205," +
					"\"public\": " +
				"\"ZG5zLXRlc3QudGhlY2RuLmV4YW1wbGUuY29tLiBJTiBETlNLRVkgMjU2IDMgOCBBd0VBQVpyMFJMdm1ubGNGK1IvMG1ESXJ2T3dZUDdJbVozaER0Tzdpc3NSQ0ZUdlBMNmhOK0dBU3ZvY3NyWXQxeTBITWt4eXN1ZnRDT25vZUh0T25QeUJaR1pWSWM2eVJFaHFLUVJlbnJ6QzBFRmVYWXhiMy9QOGFMV0pVdXBXOVdINXpRTlBEeThGVjJtUzBxSjJCNzRWYmowNnBLQlNFME1OdHQzLzZZUlJlZTlTeAo=\"," +
					"\"private\": " +
				"\"UHJpdmF0ZS1rZXktZm9ybWF0OiB2MS4yCkFsZ29yaXRobTogOCAoUlNBU0hBMjU2KQpNb2R1bHVzOiBtdlJFdSthZVZ3WDVIL1NZTWl1ODdCZy9zaVpuZUVPMDd1S3l4RUlWTzg4dnFFMzRZQksraHl5dGkzWExRY3lUSEt5NSswSTZlaDRlMDZjL0lGa1psVWh6ckpFU0dvcEJGNmV2TUxRUVY1ZGpGdmY4L3hvdFlsUzZsYjFZZm5OQTA4UEx3VlhhWkxTb25ZSHZoVnVQVHFrb0ZJVFF3MjIzZi9waEZGNTcxTEU9ClB1YmxpY0V4cG9uZW50OiBBUUFCClByaXZhdGVFeHBvbmVudDogSmtxRWpiWmduSHFpWkc0cUNnUGE3TERWVksyKzFlNU5VTmIrZkJja2JpSTEwYTVxMlRyb2tEalBMZTVPNnhTbHFlbFpFQ2ora0Z6UEcxaHg5Z2x1azUyWVlMUGJiWkFNMW1WVlc5WldRQkdDM2FzSnRRVWIxbDhLSHpldlJqN0lpNlJlcDVUeldMUnpoWGtncGdLSTV1UVBZY2xhVDJOdGk3M2kyb011SDIwPQpQcmltZTE6IDF5RERoZURDZVdORW5BSnd1ZU0xbThNdzIrb1FmaEk4Z1hwbWMvckZRdHVGMU5RS0Q0a3VEdzIvakxHakJ6RUFvNHYwZjhzQ3E5T2hJaWpvVDVSazl3PT0KUHJpbWUyOiB1R1RRdW5KZWFiS05lWW90TWUxc044M2hMZVlSNVFzY2pGalRCYXZZSzZnbG9FWWpEaSs5SzYzbWQ4YmlXN2hlU2N4bk5FejhKbFMvZ2ZNR1F5Y3hsdz09CkV4cG9uZW50MTogbVZnT1p4QzJMd2EyY2lvL0poR3lOY3hsdUd4WXd6VEdrbGlvVFFXMHRKcDhCQi84NStRVnc3OCtDZERaYjVmYlo3aXNXS2RoeVE4NkxYcFJWZUJtTXc9PQpFeHBvbmVudDI6IGE1U0dJd0Z2REFQY2ZyaWJQYkhqblh0RWtWN1Z1ZWdOcytSdTJiUTAzdU92Y0I3N2ZOOWxZd0tHb0FNdE5ZNFBsTWJvdjU3YXpoSkwyU2xNMGdrZjZRPT0KQ29lZmZpY2llbnQ6IHhDZmtZRC9UT0NvbnNoWWVxZE9yUFA4Ri93a1BaUk41S0Myc3N1U011QjlLTkNJY0VTckVTUW01K2tmbDZHUThwcGxJcVRMU2FUZWlLenFQRFU5ZGhnPT0K\"" +
				"}" +
			"]";
		final String fedZskStr =
			"[" +
				"{" +
					"\"inceptionDate\": 1545136505," +
					"\"expirationDate\": 2547728505, " +
					"\"name\": \"federation-test.thecdn.example.com.\"," +
					"\"ttl\": \"60\", " +
					"\"status\": \"new\"," +
					"\"effectiveDate\": 1543408205," +
					"\"public\": " +
				"\"ZmVkZXJhdGlvbi10ZXN0LnRoZWNkbi5leGFtcGxlLmNvbS4gSU4gRE5TS0VZIDI1NiAzIDggQXdFQUFZMFBPcVRuTjU3MU9pZTNDbTQ2aHRQU1J5NkI2dElPQUtCTExxRTMvY3ZzVFA1YzRBSUxhQ3VYWVIvd1FIUjY0MzhjWW9PWWUyM2NNeXZ5cWFZaU9vMFBCMmVpK3pER0RPZ211TmdIc3VZbXlnU1NzelY4czJSczRTNHZONjdBVjBFVVNiWCttUEFZT1IvMTJIRFRDQy9zT0hRT3IwclllREhHbCswaFBROWYK\"," +
					"\"private\":" +
				"\"UHJpdmF0ZS1rZXktZm9ybWF0OiB2MS4yCkFsZ29yaXRobTogOCAoUlNBU0hBMjU2KQpNb2R1bHVzOiBqUTg2cE9jM252VTZKN2NLYmpxRzA5SkhMb0hxMGc0QW9Fc3VvVGY5eSt4TS9semdBZ3RvSzVkaEgvQkFkSHJqZnh4aWc1aDdiZHd6Sy9LcHBpSTZqUThIWjZMN01NWU02Q2E0MkFleTVpYktCSkt6Tlh5elpHemhMaTgzcnNCWFFSUkp0ZjZZOEJnNUgvWFljTk1JTCt3NGRBNnZTdGg0TWNhWDdTRTlEMTg9ClB1YmxpY0V4cG9uZW50OiBBUUFCClByaXZhdGVFeHBvbmVudDogSmdjcEJEUGhadFV0ckc5SVBKZENxZkJTaUZNMS94TVBVQ2QwbHJvRmplaFNpWEI0WTVTM3JLak80bEZlendnaU5LNXVVSlBYRXJMK2lLYU8zZDcwY1pFR2QrUkZGWDJFUjl2VzlabmxsYlJkSDB6b3lHQkJoTk8wblIrbFVrTUhsdFBQblZjV0Fvdkt6bHZQcFMvQmVVT3R4bllZZVdUWnJzb2pBbjFKdWdFPQpQcmltZTE6IDlQK28zQkgyRTVDYWI3Z2RiZjIva1VYOFFBZ1lYcVFDb3NhbzAyN2EwOTlUdWhYc0UwaTZ3ajZ0V09ySWpmVjFGbnROTzBSZWV1emdESkszSGtIOEh3PT0KUHJpbWUyOiBrMlRCSFJTT2FKOGR5Q3p1RGZ1aUh6MmVXTEhnVGZhNzlFYXNQcGlqSGpvN3FZTGlkRWI0S2diUnpNWEVqOXUyaDRKM252c3dHVmRXRTNzMFJNS0V3UT09CkV4cG9uZW50MTogZXd6OUxxc0d3UVRiekVqWTN5bVhVY3VveWpCR3JTSUxBTjV1Wk9ORW5TMkp5K2kreldDMkRHR1doeFpFN0tmZnl3N2ExMjJiVm5vcWZhWWl1dHZCV1E9PQpFeHBvbmVudDI6IFJ5M1QrSmd4d2FKOXZtcThON0o2YzMzTlYyWG5QWjlXMnp1NStLeTdzV0JMNmF1VWNyVEhLWHlMbXNrekNJb0JWdVdSb1F3TENXSGM1cUdMOTF5OHdRPT0KQ29lZmZpY2llbnQ6IFgzeGc5cmV1cVpjWElkL0tJL2RxQWJXT05kam5BZXEwWE5paHZJQ0djQ0N6YU5GWUo0VnZuRUx0QytZbVdHT2J6bndEd3NyblVBZXRVcloxWUlkaVZRPT0K\"" +
				"}" +
			"]";
		final String cdnKskStr =
			"[" +
				"{" +
					"\"inceptionDate\": 1545136505," +
					"\"expirationDate\": 2547728505, " +
					"\"name\": \"thecdn.example.com.\"," +
					"\"ttl\": \"60\", " +
					"\"status\": \"new\"," +
					"\"effectiveDate\": 1543408205," +
					"\"public\": " +
					"\"OyBUaGlzIGlzIGEga2V5LXNpZ25pbmcga2V5LCBrZXlpZCAyOTE4MiwgZm9yIHRoZWNkbi5leGFtcGxlLmNvbS4KOyBDcmVhdGVkOiAyMDE5MDUxOTAyNDY0NCAoU2F0IE1heSAxOCAyMDo0Njo0NCAyMDE5KQo7IFB1Ymxpc2g6IDIwMTkwNTE5MDI0NjQ0IChTYXQgTWF5IDE4IDIwOjQ2OjQ0IDIwMTkpCjsgQWN0aXZhdGU6IDIwMTkwNTE5MDI0NjQ0IChTYXQgTWF5IDE4IDIwOjQ2OjQ0IDIwMTkpCnRoZWNkbi5leGFtcGxlLmNvbS4gSU4gRE5TS0VZIDI1NyAzIDggQXdFQUFhcXc2QmV4dDRMUGZzejVRVG1QeFlMWlZmcEhkSDlXS0F4T0w4d3hxdjc3OHhpdHp5bG0geEhzRzNEcGx3eHZqUTNXQUh4aE1UWmc5bmlvYXRLRkxmRXU3NkNTQlZmeVpmTDVON3R2NGlxV3kgZTM1cXRQSldxVVpkR0ZuYS9LMGR5cVFXUTgzaUpjQUhhN0NhQ2xibkRCUGRRbFZnQ0Z2NFJoM1IgcExndWsxclRFb2R1a3pZbGtEd1lVYWFtTk1qWGdnVXhpTW1jMWQ2SGp3Y200Vk8wZmhZa2hUUzUgYSt2V1ZrWXp0d3o0VGZHRnFCaXVhWDRvcTBQbGkvUnJZZVNIcWFzcWN4QUlXQTI5WUdEYjdkWlUgc2lUZTF5alRlWlhTMDZhTUxteUhCVUNJSW5iYVlrMVVFZ1JJNFAyWEtQbFEvUUpPOWNLVE9MdEcgQk1lcU56WjUzc2M9Cg==\"," +
					"\"private\": " +
					"\"UHJpdmF0ZS1rZXktZm9ybWF0OiB2MS4zCkFsZ29yaXRobTogOCAoUlNBU0hBMjU2KQpNb2R1bHVzOiBxckRvRjdHM2dzOSt6UGxCT1kvRmd0bFYra2QwZjFZb0RFNHZ6REdxL3Z2ekdLM1BLV2JFZXdiY09tWERHK05EZFlBZkdFeE5tRDJlS2hxMG9VdDhTN3ZvSklGVi9KbDh2azN1Mi9pS3BiSjdmbXEwOGxhcFJsMFlXZHI4clIzS3BCWkR6ZUlsd0FkcnNKb0tWdWNNRTkxQ1ZXQUlXL2hHSGRHa3VDNlRXdE1TaDI2VE5pV1FQQmhScHFZMHlOZUNCVEdJeVp6VjNvZVBCeWJoVTdSK0ZpU0ZOTGxyNjlaV1JqTzNEUGhOOFlXb0dLNXBmaWlyUStXTDlHdGg1SWVwcXlwekVBaFlEYjFnWU52dDFsU3lKTjdYS05ONWxkTFRwb3d1YkljRlFJZ2lkdHBpVFZRU0JFamcvWmNvK1ZEOUFrNzF3cE00dTBZRXg2bzNObm5leHc9PQpQdWJsaWNFeHBvbmVudDogQVFBQgpQcml2YXRlRXhwb25lbnQ6IG1kaGJTRWZnekNFeSs1SkpESldlQXNMYThIc0k4R0I2TmlVZWhaL2FySG52OExWdnU3UXBzVzFNZjhJS3FnOGJWVU9HUTBNNnlOWDR3YUJTWC9LR2RFaElBdWNqMWttTkdvVnBuWkFWZnlVd2s0K2Z5YkQ4WHpRM1ozMnVNbVpncDZaOXRJcDVWZXdhVHhGMzhqM0xML2hEK21sVS8zZjEwcGlMSzRxbk83cjc2NGZmWjFBanlsd2pqOUxDME1qclJoeWhhTFJoY1hudGozM0MzK3FxTGJHRGRWakRiQ0lia2UwWC9Tdy9PMVh6Q0xIdlJnZlp5NlRTaEtiOVJFK2xSYVlZdzlGZTVETEgzdG9OS1JzOXRYbldLb01OQ0ZJbXFVT1RENWxWK2Yrcm4yck9SQk5MaXRIa0NlWTJ5UkIwWXJuTFhEaHdHTGhmaDA1YUlkblZNUT09ClByaW1lMTogMGthV0VqV2tXTU52S2pOOG4yMG9TejN0aWlCQ29xOFB1c01PZUxkTVNzQUkyTTZiSC94RkJubGZ4ZDNzZnFTV2hjcWRTdlY0SzVXQmhCeWJwdEFpb1FsM1lQNnVTVjVvNnp1alMvZHI1M0ZrRTA2ZUtFUitRRWpESTZQMkhrZlVqbDNNck5qUnJ1QUMxUlBoU0dZdXE2TzladlBjalNuemd5S2t5OWhYcThrPQpQcmltZTI6IHo4N0N1ZTFZcytyaXZsWFRKQWM3QjdkZGFUSkZNVnlJN0FYUHJ4aEZ0ZmpVb1pseHBjQ1hDL2xsSHJYY05Xakx5eTF1d080aUFsSVNSQ3hpdGlWNGZZOHpaRVBpZC9VaEIvUTZCZ0ExbTE0SG5ETXJhcWdQVkJTcXRTaFZ6Rk03OFVOK1pGZVMzTFVlcllqQVF6NW5yUVd5VVVtWkpvaDd3RTloRTFHclhnOD0KRXhwb25lbnQxOiBFRm04d1oyNk1jekFrQitBeVVUTHBVNGpjbUlmekZhZ2VuMUFXdEtsOUFvS3BoRXFyc29HOUFIc0dJNnhIUWZmVEhmODB4OVRRTkJYU2RhUG8rRDdVRnBVRmc2M3JxelFxN252Y0xERWl6S2QvWUpYZWZvWmR4WXhWa3doanlrMnRmdEZOd3VGQW53WXZFalhjN0crWDBwVUovVSthUnVoKzhodDJBdnloVUU9CkV4cG9uZW50MjogcDlsQ2tianpKOGUyUTdUQTZWM3B3UzdMaFhlMFNjMkxUdERXMG4vUmRzMDR1aHBkb0ZzeDVkc1lZVGpWV0ZLQUlXbGVCdm1SZ0x4WHdyYnpPRnFGdXkwYWZvY1Nlb0FGb1E0VWU5cFpjbGY5MzUyNUdOb01INGJkNTV0ZnliMEZNcmVvZEZZRDZyOWt1eGcwNjF1UmxFQ0Fxb1crN1UvYVhSZ0F1Z0VDWU9NPQpDb2VmZmljaWVudDogWGZUQlJ0Mjdnc1I0Y2hZcGdFRG96SmppaEdxbjlUSzBSK2JDd0FDMjdQZEhCd2JjekNBYkFNZ0xValorUDVINGdqQk9Fc1VwY3l3d1Iwb1dyYXczMnlibkFETGd2K2dHU0RLM2U5T0tQaGt6Q1Zmb3l4dU9hbmVVaDBoc3liVGJ2NFVtVEV6WU5LMGxZajUvYktKV2FEZ3JOeHlsQlZDdzdUaThCcEx6TzlVPQpDcmVhdGVkOiAyMDE5MDUxOTAyNDY0NApQdWJsaXNoOiAyMDE5MDUxOTAyNDY0NApBY3RpdmF0ZTogMjAxOTA1MTkwMjQ2NDQK\"" +
				"}" +
			"]";
		final String httpsKskStr =
				"[" +
					"{" +
						"\"inceptionDate\": 1545136505," +
						"\"expirationDate\": 2547728505, " +
						"\"name\": \"https-only-test.thecdn.example.com.\"," +
						"\"ttl\": \"60\", " +
						"\"status\": \"new\"," +
						"\"effectiveDate\": 1543408205," +
						"\"public\": " +
						"\"OyBUaGlzIGlzIGEga2V5LXNpZ25pbmcga2V5LCBrZXlpZCAxNzg3MiwgZm9yIGh0dHBzLW9ubHktdGVzdC50aGVjZG4uZXhhbXBsZS5jb20uCjsgQ3JlYXRlZDogMjAxOTA1MTgyMjIzMzUgKFNhdCBNYXkgMTggMTY6MjM6MzUgMjAxOSkKOyBQdWJsaXNoOiAyMDE5MDUxODIyMjMzNSAoU2F0IE1heSAxOCAxNjoyMzozNSAyMDE5KQo7IEFjdGl2YXRlOiAyMDE5MDUxODIyMjMzNSAoU2F0IE1heSAxOCAxNjoyMzozNSAyMDE5KQpodHRwcy1vbmx5LXRlc3QudGhlY2RuLmV4YW1wbGUuY29tLiBJTiBETlNLRVkgMjU3IDMgOCBBd0VBQWI5dldyZlBXa3JmN3QySFVaM0RQNmFoZ1NIQkxZQkhnUzRYZ2x0cm0rU3NyaGg1TVpwKyB4dHFXQ2JZLzhxUjZWcEhFMyt0MWFkZzRUNnltZXlTU2NFc2huZFJFcjQ5Mm42bUpOQ2dsMlZnZSA1VVN1Tnc2T3Z3WE93eG1MYWtPdkhrNk1nRTMwcVdzTjBMd01PSUNiUkx2S0JUV0FiNlRNQVF5biBwenRlYUZIcGRnZmRBQWJRZWpIdml5YWI5cmE0WTJoMFo3TFhManhOSjlmNDB0Und4em5POUQ2aSA0NllGVW93SDR0VHpwSFliNTJ5d2QyZzBFb3RoelJhOStsRXBTeTZwU0RHcmNRQnNEYVNFZXZrbSBIWE9GMHlMcnVvc1RXdEpkRHF1SktuV2J3MTdjOVhseEp3UElKL2VrV09oZWJCQXJQTU1WZnd1YyA3SEd5WHZLeEZqcz0K\"," +
						"\"private\": " +
						"\"UHJpdmF0ZS1rZXktZm9ybWF0OiB2MS4zCkFsZ29yaXRobTogOCAoUlNBU0hBMjU2KQpNb2R1bHVzOiB2MjlhdDg5YVN0L3UzWWRSbmNNL3BxR0JJY0V0Z0VlQkxoZUNXMnViNUt5dUdIa3htbjdHMnBZSnRqL3lwSHBXa2NUZjYzVnAyRGhQcktaN0pKSndTeUdkMUVTdmozYWZxWWswS0NYWldCN2xSSzQzRG82L0JjN0RHWXRxUTY4ZVRveUFUZlNwYXczUXZBdzRnSnRFdThvRk5ZQnZwTXdCREtlbk8xNW9VZWwyQjkwQUJ0QjZNZStMSnB2MnRyaGphSFJuc3RjdVBFMG4xL2pTMUhESE9jNzBQcUxqcGdWU2pBZmkxUE9rZGh2bmJMQjNhRFFTaTJITkZyMzZVU2xMTHFsSU1hdHhBR3dOcElSNitTWWRjNFhUSXV1Nml4TmEwbDBPcTRrcWRadkRYdHoxZVhFbkE4Z245NlJZNkY1c0VDczh3eFYvQzV6c2NiSmU4ckVXT3c9PQpQdWJsaWNFeHBvbmVudDogQVFBQgpQcml2YXRlRXhwb25lbnQ6IE44MmVCRGJOZTBZTHUwZlc0c1lucDhzc2VVcDJtUTQrK2RDZ2owV3ZDOW5LWmhmdC9iczIvRUVBVThBUVd5SE9XbStwVmxuRG9PUEpWZXF4dXRkMUpIR0lNSGhWTk55L2Jnd3d5QU5BZUErSmhadkRNTnNyaytYUnVZQ0tXWENTeFJMdjA4bWVHVGJOd2dOTjlTOU51ZkFKMUs2NzNLNGJJRFUrNm05NnVXVnphU0ovaDgreXdyQk5PV0xsZTNTQzV2TEFYcnhseWliVDlDSm5ueExGVGxINGtCanY3NHZoSHlpcUhyaVhGZ3NjOFU2eVVIMEY2cHptTTlEUEV3TG4vYjZTREtWYzhFZ3o2V2JmZzZVTUF5VjFJbVN0aWU1L2NRYzN2d2tya3JOMU8zTDZZZVhjM3d0QW1MQ3U0Q3hFYnNZLy9ZcXFFbDZIRk8wd1M2SEZBUT09ClByaW1lMTogNkxHMEk0L2RDRlhMZ3o4Wmd1VnozNzV3UlRRcGRyVWtNNXo3MjkxTmVUZXJGMGxzTkRZWDBXejh3K0dmRmdpZ2RxWmZmdHJNYmYvWVpCMnkyc2Zjc1RsTnFRYmw4VGV6QVJhWjJ2U09mRFJ0TGI0T2JRdkxzcUpnRE5tUnhXQTZGZUtuVUorYk5TY1h6YzZwOTZVTjRYWkwxTmNRT2NMMmc4ZDRvYlJiaEVjPQpQcmltZTI6IDBwdkRIZTgrSVcwSnFRSGhQQmdhRllWRzdsL1k5MVUvdE5WOVBHbTBuT1VaRFBUNnByOVFWaGVmWFhVNVA4SmJkRURNMlVLV0ZUaFBTNWIzVy9LTHRCdHZJcWxpdm9kNWgyTWE5aEZsUDh2dzNyL2E4ZmI0b1ZVMG9vcFh2Qy83RmVwTnl5b0ZJV1R4RzRBNVVGQStlRHFJbHB5TVBuWnE1RlNlM2x1dUhHMD0KRXhwb25lbnQxOiA1YWE4U255cGdKaHNDbFEwTVdPVFFMY0t4c0g4U2hQc2JxUDRUYjNUd0ZhWW5KcnlGM1ZyZkYwNytYYXJNMnZBTWxsdzFobkt1S1ZRUXo2c1RoQUNWMFpleHZydjVXazdXVStjK09OejNGRkJqMnVMZ1VPcS9kb1RRWnRZcXB1VnVCUEJYV2lvSFlVL2tQYnQrR01GbUFiUVFIY2dwR0V1T2xDYlZieFN0ZkU9CkV4cG9uZW50MjogdGUxZmF1aFRYMFIxWjh6NzU1RmFWdVMrRlFRdXc5aWNJM1dYclN3U25NVTZFbnM4V2ZaQlMxMDBpT0xPQlVtNi9uMUxkeEdSMjlxOGhLdHdHYmsyL09vRjRvYzNpU1kxME1ISGRIQXFhaVdkZUkxNmNESExMSElSK2FaUGkzeFhCT05WTi82Z1YreCthaWNsV3o4MTkxMTR4OEdMVkJtdTFIWlVsZmZVT3pFPQpDb2VmZmljaWVudDogdjUwbmtGM1NmTnd3WE5DcFF4ZVIzZ0MzVU1FRnc4SkJrcDZwZjQyQlZJV05EVm1TMEIrQXYxL1EwdnFnZC80N0pWaFp6MjlvTHBCdWNKdDF0cUtrdkpqekI5dmxweDRWUHVCY0JSNkpTMzJrRkNDQUtLcjB1S1JnY3U4Y0lVcUtuS3BPb1JIR0lYc1BBRHVEWVRtWGZXYnlGR1VXR25mMzVBbW9jMmFscTlVPQpDcmVhdGVkOiAyMDE5MDUxODIyMjMzNQpQdWJsaXNoOiAyMDE5MDUxODIyMjMzNQpBY3RpdmF0ZTogMjAxOTA1MTgyMjIzMzUK\"" +
					"}" +
				"]";
		final String dnstKskStr =
			"[" +
				"{" +
					"\"inceptionDate\": 1545136505," +
					"\"expirationDate\": 2547728505, " +
					"\"name\": \"dns-test.thecdn.example.com.\"," +
					"\"ttl\": \"60\", " +
					"\"status\": \"new\"," +
					"\"effectiveDate\": 1543408205," +
					"\"public\": " +
				"\"ZG5zLXRlc3QudGhlY2RuLmV4YW1wbGUuY29tLiBJTiBETlNLRVkgMjU3IDMgOCBBd0VBQVozU09rTjZ1bnVxYlM5ZGtDcnE4VFQyT1JTcmNrOHE3bVZDUEhtMmxYKzdBTHU3OURsOE9nVFEvTkxTd09iNk0wNmo3QW0wT0ROZElJVllqeGFuRXNqRWZ3c1RUUFg0MDhHc0NPa1BOeHNoclZMU0ZXUEJ3dXF6SW1VVElOT0MyVHByckMwNkswRzJCNFVhbG9CTElZTEsxOTRwT2VHK1FVQ0p4ZkJERWZjVEh0ZWdLOHlvc29MamNYZTM5L3k5RW5kV2YxV2JWZTE0RTNhTmdqeDJlcjlwNnl4MkJJZHVHOGJ6Y2dGS3AzbHFEdjFrZE9tU2ZTSXV0Rm50TzhPQkJ1M25DYisvWWtpNlE3TkROd1pTaHIxazdHTXFmOTZqbEZVNUhGYUdiN2xlNklYVXh3YWtzUWdrZFQ3THhYVUJoVjhlWWxNVFhsV3NsTzRQcXJoVXA2Yz0K\"," +
					"\"private\": " +
				"\"UHJpdmF0ZS1rZXktZm9ybWF0OiB2MS4yCkFsZ29yaXRobTogOCAoUlNBU0hBMjU2KQpNb2R1bHVzOiBuZEk2UTNxNmU2cHRMMTJRS3VyeE5QWTVGS3R5VHlydVpVSThlYmFWZjdzQXU3djBPWHc2Qk5EODB0TEE1dm96VHFQc0NiUTRNMTBnaFZpUEZxY1N5TVIvQ3hOTTlmalR3YXdJNlE4M0d5R3RVdElWWThIQzZyTWlaUk1nMDRMWk9tdXNMVG9yUWJZSGhScVdnRXNoZ3NyWDNpazU0YjVCUUluRjhFTVI5eE1lMTZBcnpLaXlndU54ZDdmMy9MMFNkMVovVlp0VjdYZ1RkbzJDUEhaNnYybnJMSFlFaDI0Ynh2TnlBVXFuZVdvTy9XUjA2Wko5SWk2MFdlMDd3NEVHN2VjSnY3OWlTTHBEczBNM0JsS0d2V1RzWXlwLzNxT1VWVGtjVm9adnVWN29oZFRIQnFTeENDUjFQc3ZGZFFHRlh4NWlVeE5lVmF5VTdnK3F1RlNucHc9PQpQdWJsaWNFeHBvbmVudDogQVFBQgpQcml2YXRlRXhwb25lbnQ6IGh4RzZUYkJHMDdvTFVpTllWSExZMXdQMzNFblRQaEEzRWJCN2c0dVJMVTFGbG1hSTRYNEJSY2Y2NlEvNGluWU4zVHNMczA1clh3SlA1Ky9nSG5vRTZKRExUaFpKb3FZL3pSeElUL1oycWlETGJ2dGYxUTJxblNXTXhVWjJySzdxN1VYamlKMmxFY3NSYW9oVDBCNzg0aXhxVGJlbzB4djZTcHJmTGY2bzdIUXh6bVlGZVMxaGdRbWhuZm1rTnp6dElpMEtjcDFWMG1Pa2NBcmhteFpkWTgvUHJsalpCL3RSMTZQdU1KbG16dG1aUE95Wm1Gbnk4c3RiYWt4VkppL0s3Uit0aHU1NHVTYU1vamxuV3FuVnh0U0QyM1VQaFQzZDFFVERyVGd1UlBQbXB4YU02aWxKYUdERkpPL3ZzenhEaE5YT1ZGU3E3M1FoVUpLeDEySncwUT09ClByaW1lMTogOGZieFJHY0ordEV1SDN5RHI5cTZqVkdySzF6K09ENlA5d2JOeEJLR0NUNndCM3RXNmNzUkpQOThIeEovSUEvenRtL2MrTnRHRVBmdDZLVDZLNjVCMkNKY3Z3WmZmcDZZOEpSekF3dzZEWWV0MHFpeGttZWZqYUVKZWtWSUxGUTZSdU9nbllCdUVNTFpaby93N3M4K0VXbUdyK0xMR3RyR3BMU2VNZmkyOHY4PQpQcmltZTI6IHB2bkx6TjU2NFpCaWRycXZFaVg1dC91Q1dSYlRNdzNQNEVnNThoWUZkR1FpbmVIUm53d2tzRUlleG5Id2VETExtWWZXcEsxZkNLck9aeUhiUEZ5OFdhY3ZEUkt5Q3AzeFRTcEhRSmk0L2xXcTdlVHFQaEpJV2YvYVVkeDJ1RW1kbzNwREFVQlZ2eDhZUlBkUjR6Szk3QnhkSHFHNTJ3WVltZzNhOTFPTzAxaz0KRXhwb25lbnQxOiBQdUxvZDllejMwMUlpSVJyRVdSdXdkWHMvOXN1YzEzSE92TzR2UEgzaGlXVnlJd0UzY1NhVXh4WG5SZklsSU93MnNTZUVNdWtuVHBpeWVrKzMrVnRWWWd3eExFYVZxVlBxSTljaVBrL2lVNnZIYVljYUttbjdUNWlZVFhxZVNMMjluK291ZWFzTkl6L3hjazVYRWZlb05YbFhJYzhOR0dSNlRMTVByNmVoZTg9CkV4cG9uZW50MjogZzQ4RkdDR2l4OTR1OWtVWWMwQWdoT2xSUmtoSmwwd21vUnZITEFwVnVlSzdzNUdjeTZlUnNKNG9DVXIwb0gvRkV1NklHNi9OMU5KZlZickROY2dMVHNmK3Rsb29sVnprSmx4TlQ0UUZIYjc1c2Y1TzRTRWVpR3FoNVNYREZHaE1IK1hRclVlM1I2S0VTTEprZnBJWU9kUVBPbmRLTEZ1ZFBxUDBCako3c2VFPQpDb2VmZmljaWVudDogald3L2phSmNJeVRtclhNc3J2L1NucmRjUHEzZmlpVkc1K3BZYmVYN1g1MFBMMVJtWE5VM2JVM09SdmlNZE9mbjdlR2dXaWkrNnFZQVdaZDBqbDJ6em5na1NMUnkzZytQNWpya0VyelNZVStzaTM1RklJUURObnU2YUlHUXJXS25oa29ESEhORzM0MXJFSytHRHZXbXg4dWpjTjVQTzk5ZUdYNEE4L0RaTFQ0PQo=\"" +
				"}" +
			"]";
		final String fedKskStr =
			"[" +
				"{" +
					"\"inceptionDate\": 1545136505," +
					"\"expirationDate\": 2547728505, " +
					"\"name\": \"federation-test.thecdn.example.com.\"," +
					"\"ttl\": \"60\", " +
					"\"status\": \"new\"," +
					"\"effectiveDate\": 1543408205," +
					"\"public\": " +
				"\"ZmVkZXJhdGlvbi10ZXN0LnRoZWNkbi5leGFtcGxlLmNvbS4gSU4gRE5TS0VZIDI1NyAzIDggQXdFQUFaRnVUakFqUWN6V3ZjVWtITlQvQXBOdGZIRVEwd2QwNElPNStmY2laK3VWV3ZwNjVZcm1RTWI5WXZBcnY3aUhSYzFweDRsbXVxK0VndHBQekZEN3N6VDk0cEdHWXhzcTR2RG5BWFBad1h3clBuUjJBWkZVS1owWDAwQlBXdWdXRDcvUHJuRXAxVTg1V3dhSFBsd3JiSWhlWW01L1E3aTRUY1BvT0tHVjl1SzJyQXlnbzJ5d2dhd0NBWFp2cUp4Smg2bTRwVWx4Ry9YdVUxM0NLSFVRRURJSHo5UnliT0ZHcTAzMFZqbU92UndaK21DVGFQbmFsZTZjUEhFU1hXNms1aTVBeDA5cXBwTElFRkYySy91YzlRdGZlMXpjdUxkNVRMZ0hWcVhMYVA5RkxMU1ZPbXRmVFM3SEQ5WG8wcTc2ZGFSQktjSWEzUWIvSjFoU1NIZEU3SXM9Cg==\"," +
					"\"private\":" +
				"\"UHJpdmF0ZS1rZXktZm9ybWF0OiB2MS4yCkFsZ29yaXRobTogOCAoUlNBU0hBMjU2KQpNb2R1bHVzOiBrVzVPTUNOQnpOYTl4U1FjMVA4Q2syMThjUkRUQjNUZ2c3bjU5eUpuNjVWYStucmxpdVpBeHYxaThDdS91SWRGelduSGlXYTZyNFNDMmsvTVVQdXpOUDNpa1laakd5cmk4T2NCYzluQmZDcytkSFlCa1ZRcG5SZlRRRTlhNkJZUHY4K3VjU25WVHpsYkJvYytYQ3RzaUY1aWJuOUR1TGhOdytnNG9aWDI0cmFzREtDamJMQ0JyQUlCZG0rb25FbUhxYmlsU1hFYjllNVRYY0lvZFJBUU1nZlAxSEpzNFVhclRmUldPWTY5SEJuNllKTm8rZHFWN3B3OGNSSmRicVRtTGtESFQycW1rc2dRVVhZcis1ejFDMTk3WE55NHQzbE11QWRXcGN0by8wVXN0SlU2YTE5TkxzY1AxZWpTcnZwMXBFRXB3aHJkQnY4bldGSklkMFRzaXc9PQpQdWJsaWNFeHBvbmVudDogQVFBQgpQcml2YXRlRXhwb25lbnQ6IGg0QzJXMFhPZmxRclZ5OHhxZ2U4MTY2d3Z3eUZBN0tUcWtpekxlQWg0YkEwcDZPd2tuMjlKMnRhTHhza05JUGR0dW56WUFPV3VBa0lmdTdSR1RlY0h5amJYT3BSRnpRYlpZaG5veERtcFpJSlRDdlRoQnhkOWFBSVZpaGFORnF4Nis5T3d1UE9lMVdlaVhPajErOGgzZUhMWnRjdk8wS0dPcDM1ZmgwamZ0SjdZVmU2N0tGdzJBT1lQK1dZalBFaDJpb2ZWT2NVNWorR3VYWkt3V1NHUmZzM0lRYlp5Y2Q0SHZ0YTZmbndBR2s3WGhrSkhsQ3lBeHNvMTBCNVdmZHJ5ellNZ29LV2lKOHRZUGljTVVGcVJOaHM2UElYUTI0Qk9ITjdxQmF3UTJ6QXJralB6L3B4aWRKRTc3OURHQ29zZWFxRURrdytVak93SDV5N2RzU0lDUT09ClByaW1lMTogeTVMT0ZkVXEvZjVSSFpWR3pob2ozYkF2MU5PWEJQYjk3TXZmYmF5eFNWVnJiaGJzQXFMY05XTUFCdC9mRHpMQ1FZWjBneW4zR1cxWXBENWlUOG9SZzM2V3AwZ3gwWGs4dU50NEpxalpNWE5ySDZtVWQ2T0Jydm81dzl3d0lSQWY1VDh6UWJPVkE0N1I5RHJHTjk3RlJqYkNBVW1pWFpmN3E4ZUJuczF1U2tVPQpQcmltZTI6IHR1SkdnM1lvTG0yeE9YVHQ4TEFsdWt0akJJOGhXMncxTUlXSDFuRDBHQ3ZMY2ZPSUVtQ0ZGMUt1TkluaDhqcXZ5TUNjekdyU3k0bGF5WHg0OUZDV0dFT3ZuQ1JBLzcrRGo0RVJrS1pNdkZPNmxHYWV3Mm0vSkZZVGZEblNweEpXSkJtSW9pVmJkYVhsQ2hjUjM2SkNEYUFKT1J4b29aVS9pSHN1MWd2U3NJOD0KRXhwb25lbnQxOiBWU2hSSzFManpDSlJubFZ1ckJMRlJCeEt0ZlhaSzh1Q2gwYjFiUFNicVBpaG13amRxM0NqTzNYeGNlNithYVlyR3F2N0cwODN2WncvUTEyUlZKMUwzRHpkR3BjWnQrM0dWL0grL2ZVTi9pQ3hCQ3ExSDZMM1FkSU16Z0RTNVZIUWRkNk5PNE82NXlVY2NOVVJUQmZWWUR6UnhTWWZWSldhUXM2UFMzWFdHQjA9CkV4cG9uZW50MjogTmV2emRIRlRHWlZZQ3FQU1FBUC9xN1RzaGZ5WmpqWVNYTE1TUVFUZXczMnVKM1B4YTlHdmpCZmhxelg0TzQ1WUkrMitqWHIxbWZOdXBEZWlCZzc0b2tEYXQwUHRNanJLVkhadXNtS0YvNFVFWHhyK3Rva29SVk5udlZuakpVVi94bmNNMVJvRXBHUjhhb1F3emVvdVpZd0pEQ0MzTE9VdmJWTThsUG01YmpzPQpDb2VmZmljaWVudDogTkR3ditiMnlhWnRpNVFGTXRtbFBjQmpXdW5kTWd5TXp3REs0MUUzajI0QzVoMzQ4TG55UFB2aDNTMlBjYTZpanpydTdaSVRPNjdRZFpWL04rcTdrZ0ZIbTIzZ3E0V1c4dkxqQnhWMEZFcEJCWWRZK0pvOUdjYUJxK29LekUyblFQdklMZkFka1VsMExZWkoxa1FpMTVweHZNbERJazBDempLNDNNMVphT2wwPQo=\"" +
				"}" +
			"]";

		final String dnstKpdStr =
			"{" +
				"\"response\": {" +
					"\"thecdn\": {" +
						"\"zsk\": " + cdnZskStr + "," +
						"\"ksk\": " + cdnKskStr +
					"}," +
					"\"https-only-test\": {" +
						"\"zsk\": " + httpsZskStr + "," +
						"\"ksk\": " + httpsKskStr +
					"}," +
					"\"dns-test\": {" +
						"\"zsk\": " + dnstZskStr + "," +
						"\"ksk\": " + dnstKskStr +
					"}" +
				"}" +
			"}";

		final String dnsFedKpdStr =
			"{" +
				"\"response\": {" +
					"\"thecdn\": {" +
						"\"zsk\": " + cdnZskStr + "," +
						"\"ksk\": " + cdnKskStr +
					"}," +
					"\"https-only-test\": {" +
						"\"zsk\": " + httpsZskStr + "," +
						"\"ksk\": " + httpsKskStr +
					"}," +
					"\"dns-test\": {" +
						"\"zsk\": " + dnstZskStr + "," +
						"\"ksk\": " + dnstKskStr +
					"}," +
					"\"federation-test\": {" +
						"\"zsk\": " + fedZskStr + "," +
						"\"ksk\": " + fedKskStr +
					"}" +
				"}" +
			"}";

		try{
			final ObjectMapper mapper = new ObjectMapper();
			if (returnKey == KeyProfile.TWO){
				return mapper.readTree(dnsFedKpdStr);
			} else {
				return mapper.readTree(dnstKpdStr);
			}
		} catch (IOException ioe){
			fail(ioe.getMessage());
			return null;
		}
	}

	@Override
	protected boolean isDnssecEnabled() {return true;}
}

@RunWith(PowerMockRunner.class)
@PrepareForTest({ConfigHandler.class, CacheRegister.class, ZoneManager.class,
		TrafficRouterManager.class, TrafficRouter.class })
public class SignatureManagerTest {
	ZoneManager zoneManager;
	SignatureManager signatureManager;
	CacheRegister cacheRegister;
	LoadingCache<ZoneKey, Zone> dynamicZoneCache;
	LoadingCache<ZoneKey, Zone> zoneCache;
	String baseDb = null;
	JsonNode baselineJo, updateJo, modDsJo = null;
	TrafficRouterManager trafficRouterManager;
	ConfigHandler configHandler = PowerMockito.spy(new ConfigHandler());
	StatTracker statTracker = new StatTracker();
	TrafficRouter trafficRouter = null;

	public static String DNS_TEST = "dns-test.thecdn.example.com.";
	public static String FED_TEST = "federation-test.thecdn.example.com.";
	public static String HTTPS_TEST = "https-only-test.thecdn.example.com.";
	public static String TLD = "thecdn.example.com.";

	@Before
	public void before() throws Exception{
		try{
			System.setProperty("dns.zones.dir", "src/test/var/auto-zones");
			String resourcePath = "unit/DNSSecCrConfig.json";
			InputStream inputStream = getClass().getClassLoader().getResourceAsStream(resourcePath);
			if (inputStream == null){
				fail("Could not find file '" + resourcePath + "' needed for test from the current classpath as a resource!");
			}
			baseDb = IOUtils.toString(inputStream);

			resourcePath = "unit/DNSSecPlusCrConfig.json";
			inputStream = getClass().getClassLoader().getResourceAsStream(resourcePath);
			if (inputStream == null){
				fail("Could not find file '" + resourcePath + "' needed for test from the current classpath as a resource!");
			}
			String newDsSnap = IOUtils.toString(inputStream);
			resourcePath = "unit/DNSSecModDsCrConfig.json";
			inputStream = getClass().getClassLoader().getResourceAsStream(resourcePath);
			if (inputStream == null){
				fail("Could not find file '" + resourcePath + "' needed for test from the current classpath as a resource!");
			}
			String modDsSnap = IOUtils.toString(inputStream);
			final ObjectMapper mapper = new ObjectMapper();
			assertThat(baseDb, notNullValue());
			assertThat(newDsSnap, notNullValue());
			assertThat(modDsSnap, notNullValue());

			baselineJo = mapper.readTree(baseDb);
			updateJo = mapper.readTree(newDsSnap);
			modDsJo = mapper.readTree(modDsSnap);
			assertThat(baselineJo, notNullValue());

			GeolocationService geolocationService = new MaxmindGeolocationService();
			AnonymousIpDatabaseService anonymousIpDatabaseService = new AnonymousIpDatabaseService();
			FederationRegistry federationRegistry = new FederationRegistry();
			TrafficOpsUtils trafficOpsUtils = new TrafficOpsUtils();
			trafficRouterManager = new TrafficRouterManager();
			trafficRouterManager.setAnonymousIpService(anonymousIpDatabaseService);
			trafficRouterManager.setGeolocationService(geolocationService);
			trafficRouterManager.setGeolocationService6(geolocationService);
			trafficRouterManager.setFederationRegistry(federationRegistry);
			trafficRouterManager.setTrafficOpsUtils(trafficOpsUtils);
			trafficRouterManager.setNameServer(new NameServer());
			configHandler.setTrafficRouterManager(trafficRouterManager);
			configHandler.setStatTracker(statTracker);
			configHandler.setFederationsWatcher(new FederationsWatcher());
			configHandler.setSteeringWatcher(new SteeringWatcher());
			SnapshotEventsProcessor snapshotEventsProcessor = SnapshotEventsProcessor.diffCrConfigs(baselineJo, null);
			Map<String, DeliveryService> deliveryServiceMap = snapshotEventsProcessor.getChangeEvents();
			ZoneManager.setZoneDirectory(new File("src/test/resources/unit/sigmantest"));
			cacheRegister = PowerMockito.spy(new CacheRegister());
			cacheRegister = fillCacheRegister(cacheRegister, deliveryServiceMap, null, baselineJo);
			trafficRouter = PowerMockito.mock(TrafficRouter.class);
			when(trafficRouter, "getCacheRegister").thenReturn(cacheRegister);
			zoneManager = PowerMockito.spy(new ZoneManager(trafficRouter, statTracker));

			// These three lines only allow the SignatureManager to be created once. Therefore each unit test is
			// interacting with the same SignatureManager instance as the previous test.
			signatureManager = new SigManagerForTesting(zoneManager, cacheRegister, trafficOpsUtils, trafficRouterManager);
			whenNew(SignatureManager.class).withArguments(zoneManager, cacheRegister,
					trafficOpsUtils, trafficRouterManager).thenReturn(signatureManager);
			///////////////
			SigManagerForTesting.returnKey = SigManagerForTesting.KeyProfile.ONE;
			Whitebox.invokeMethod(zoneManager, "initDnsRoutingNames", cacheRegister);
			Whitebox.invokeMethod(zoneManager, "initTopLevelDomain", cacheRegister);
			Whitebox.invokeMethod(zoneManager, "initSignatureManager", cacheRegister, trafficOpsUtils, trafficRouterManager);
			Whitebox.invokeMethod(zoneManager, "initZoneCache", trafficRouter);
			zoneCache = Whitebox.getInternalState(ZoneManager.class, "zoneCache");
			dynamicZoneCache = Whitebox.getInternalState(ZoneManager.class, "dynamicZoneCache");
		} catch (Exception ex){
			ex.printStackTrace();
			fail(ex.getMessage());
		}
	}

	private CacheRegister fillCacheRegister(CacheRegister cacheRegister,
	                                        Map<String, DeliveryService> deliveryServiceMap,
	                                        SnapshotEventsProcessor snapshotEventsProcessor, JsonNode snapshotJo) throws
			Exception{
		final JsonNode config = JsonUtils.getJsonNode(snapshotJo, ConfigHandler.CONFIG_KEY);
		final JsonNode stats = JsonUtils.getJsonNode(snapshotJo, "stats");
		cacheRegister.setTrafficRouters(JsonUtils.getJsonNode(snapshotJo, "contentRouters"));
		cacheRegister.setConfig(config);
		cacheRegister.setStats(stats);
		Whitebox.invokeMethod(configHandler, "parseCertificatesConfig", config);
		if (snapshotEventsProcessor == null){
			Whitebox.invokeMethod(configHandler, "parseDeliveryServiceMatchSets", deliveryServiceMap, cacheRegister);
			Whitebox.invokeMethod(configHandler, "parseLocationConfig", JsonUtils
					.getJsonNode(snapshotJo, "edgeLocations"), cacheRegister);
			Whitebox.invokeMethod(configHandler, "parseCacheConfig", JsonUtils.getJsonNode(snapshotJo,
					ConfigHandler.CONTENT_SERVERS_KEY), cacheRegister);
		}else{
			Whitebox.invokeMethod(configHandler, "parseDeliveryServiceMatchSets", snapshotEventsProcessor, cacheRegister);
			Whitebox.invokeMethod(configHandler, "parseLocationConfig", JsonUtils
					.getJsonNode(snapshotJo, "edgeLocations"), cacheRegister);
			Whitebox.invokeMethod(configHandler, "parseCacheConfig", snapshotEventsProcessor,
					JsonUtils.getJsonNode(snapshotJo, ConfigHandler.CONTENT_SERVERS_KEY), cacheRegister);
		}
		Whitebox.invokeMethod(configHandler, "parseMonitorConfig", JsonUtils.getJsonNode(snapshotJo, "monitors"));
		return cacheRegister;
	}

	@Test
	public void verifyInitialState(){
		Map<String, List<DnsSecKeyPair>> keyMap = signatureManager.getKeyMap();
		assertThat("keyMap should not be null.", keyMap, notNullValue());
		assertThat("Expected to only find a keys for [" + TLD + ", " + DNS_TEST + ", " + HTTPS_TEST +
						"] but found these keys - " + keyMap.keySet(),
				keyMap.keySet(), containsInAnyOrder(TLD, DNS_TEST, HTTPS_TEST));
	}

	@Test
	public void snapNewDNSDsNoNewKeys(){
		try{
			SnapshotEventsProcessor snapshotEventsProcessor = SnapshotEventsProcessor.diffCrConfigs(updateJo,
					baselineJo);
			cacheRegister = fillCacheRegister(cacheRegister, null, snapshotEventsProcessor, updateJo);
			when(trafficRouter, "getCacheRegister").thenReturn(cacheRegister);
			when(zoneManager, "getTrafficRouter").thenReturn(trafficRouter);
			signatureManager.refreshKeyMap();
		} catch (Exception e){
			fail(e.getMessage());
		}

		verify(zoneManager, never()).updateZoneCache(anyList());
		Map<String, List<DnsSecKeyPair>> keyMap = signatureManager.getKeyMap();
		assertThat("keyMap should not be null.", keyMap, notNullValue());
		assertThat("Expected to only find a keys for [" + TLD + ", " + DNS_TEST + ", " + HTTPS_TEST +
						"] but found these keys - " + keyMap.keySet(),
				keyMap.keySet(), containsInAnyOrder(TLD, DNS_TEST, HTTPS_TEST));
	}

	@Test
	public void snapNewDNSDsAndNewKeys(){
		Map<ZoneKey, Zone> cacheOut = null;
		try{
			SnapshotEventsProcessor snapshotEventsProcessor = SnapshotEventsProcessor.diffCrConfigs(updateJo,
					baselineJo);
			cacheRegister = fillCacheRegister(cacheRegister, null, snapshotEventsProcessor, updateJo);
			SigManagerForTesting.returnKey = SigManagerForTesting.KeyProfile.TWO;
			when(trafficRouter, "getCacheRegister").thenReturn(cacheRegister);
			when(zoneManager, "getTrafficRouter").thenReturn(trafficRouter);
			signatureManager.refreshKeyMap();
			zoneCache = Whitebox.getInternalState(ZoneManager.class, "zoneCache");
			cacheOut = zoneCache.asMap();
		} catch (Exception e){
			fail(e.getMessage());
		}

		verify(zoneManager, atLeastOnce()).updateZoneCache(anyList());
		Map<String, List<DnsSecKeyPair>> keyMap = signatureManager.getKeyMap();
		assertThat("keyMap should not be null.", keyMap, notNullValue());
		assertThat("Expected to find a keys for [" + TLD + ", " + DNS_TEST + ", " + HTTPS_TEST + ", " + FED_TEST +
						"] but found these keys - " + keyMap.keySet(),
				keyMap.keySet(), containsInAnyOrder(TLD, DNS_TEST, HTTPS_TEST, FED_TEST));

		final StringBuilder keysStr = new StringBuilder("");
		cacheOut.keySet().forEach(zoneKey -> keysStr.append(zoneKey.getName() + " "));
		assertThat("Expected a zone for 'dns-test.thecdn.example.com.' in :" + keysStr.toString(),
				keysStr.toString(), containsString(DNS_TEST));
		assertThat("Expected a zone for 'federation-test.thecdn.example.com.' in :" + keysStr.toString(),
				keysStr.toString(), containsString(FED_TEST));
		Zone dnsTestZone = null;
		Zone fedTestZone = null;
		Zone tldTestZone = null;
		for (ZoneKey key : cacheOut.keySet()){
			if (key.getName().toString().equals(DNS_TEST)){
				dnsTestZone = cacheOut.get(key);
			}else if (key.getName().toString().equals(FED_TEST)){
				fedTestZone = cacheOut.get(key);
			}else if (key.getName().toString().equals(TLD)){
				tldTestZone = cacheOut.get(key);
			}

			if (!(dnsTestZone == null || fedTestZone == null || tldTestZone == null)){
				break;
			}
		}

		assertThat("Check DNSSec Key for " + DNS_TEST, hasNSec(dnsTestZone, DNS_TEST, Type.DNSKEY),
				notNullValue());
		assertThat("Check NSEC Key for " + DNS_TEST, hasNSec(dnsTestZone, DNS_TEST, Type.NSEC),
				notNullValue());
		assertThat("Check DNSSec Key for " + FED_TEST, hasNSec(fedTestZone, FED_TEST, Type.DNSKEY), notNullValue());
		assertThat("Check NSEC Key for " + FED_TEST, hasNSec(fedTestZone, FED_TEST, Type.NSEC), notNullValue());
		assertThat("Check DS record for " + DNS_TEST, hasNSec(tldTestZone, DNS_TEST, Type.DS), notNullValue());
		assertThat("Check DS record for " + FED_TEST, hasNSec(tldTestZone, FED_TEST, Type.DS), notNullValue());
	}

	@Test
	public void snapModDsAndSameKeys(){
		Map<ZoneKey, Zone> cacheOut = null;
		try{
			SnapshotEventsProcessor snapshotEventsProcessor = SnapshotEventsProcessor.diffCrConfigs(modDsJo,
					baselineJo);
			cacheRegister = fillCacheRegister(cacheRegister, null, snapshotEventsProcessor, modDsJo);
			signatureManager.getKeyMap().remove(FED_TEST);
			SigManagerForTesting.returnKey = SigManagerForTesting.KeyProfile.TWO;
			when(trafficRouter, "getCacheRegister").thenReturn(cacheRegister);
			when(zoneManager, "getTrafficRouter").thenReturn(trafficRouter);
			signatureManager.refreshKeyMap();
			zoneCache = Whitebox.getInternalState(ZoneManager.class, "zoneCache");
			cacheOut = zoneCache.asMap();
		} catch (Exception e){
			fail(e.getMessage());
		}

		verify(zoneManager, atLeastOnce()).updateZoneCache(anyList());
		Map<String, List<DnsSecKeyPair>> keyMap = signatureManager.getKeyMap();
		assertThat("keyMap should not be null.", keyMap, notNullValue());
		assertThat("Expected to only find a keys for [" + DNS_TEST + ", " + HTTPS_TEST + ", " + FED_TEST +
						"] but found these keys - " + keyMap.keySet(),
				keyMap.keySet(), containsInAnyOrder(TLD, DNS_TEST, HTTPS_TEST, FED_TEST));

		final StringBuilder keysStr = new StringBuilder("");
		cacheOut.keySet().forEach(zoneKey -> keysStr.append(zoneKey.getName() + " "));
		assertThat("Expected a zone for 'dns-test.thecdn.example.com.' in :" + keysStr.toString(),
				keysStr.toString(), containsString(DNS_TEST));
		assertThat("Expected a zone for 'federation-test.thecdn.example.com.' in :" + keysStr.toString(),
				keysStr.toString(), containsString(FED_TEST));
		Zone dnsTestZone = null;
		Zone fedTestZone = null;
		Zone tldTestZone = null;
		for (ZoneKey key : cacheOut.keySet()){
			if (key.getName().toString().equals(DNS_TEST)){
				dnsTestZone = cacheOut.get(key);
			}else if (key.getName().toString().equals(FED_TEST)){
				fedTestZone = cacheOut.get(key);
			}else if (key.getName().toString().equals(TLD)){
				tldTestZone = cacheOut.get(key);
			}

			if (!(dnsTestZone == null || fedTestZone == null || tldTestZone == null)){
				break;
			}
		}

		assertThat("Check DNSSec Key for " + DNS_TEST, hasNSec(dnsTestZone, DNS_TEST, Type.DNSKEY),
				notNullValue());
		assertThat("Check NSEC Key for " + DNS_TEST, hasNSec(dnsTestZone, DNS_TEST, Type.NSEC),
				notNullValue());
		assertThat("Check DNSSec Key for " + FED_TEST, hasNSec(fedTestZone, FED_TEST, Type.DNSKEY), notNullValue());
		assertThat("Check NSEC Key for " + FED_TEST, hasNSec(fedTestZone, FED_TEST, Type.NSEC), notNullValue());
		assertThat("Check DS record for " + DNS_TEST, hasNSec(tldTestZone, DNS_TEST, Type.DS), notNullValue());
		assertThat("Check DS record for " + FED_TEST, hasNSec(tldTestZone, FED_TEST, Type.DS), notNullValue());
	}

	public Record hasNSec(final Zone srcZone, final String hostname, final int type){
		Iterator<RRset> rRsetIterator = srcZone.iterator();
		while (rRsetIterator.hasNext()){
			RRset recordSet = rRsetIterator.next();
			Iterator<Record> recordIterator = recordSet.rrs();
			while (recordIterator.hasNext()){
				Record record = recordIterator.next();
				if (record.getType() == type && record.rdataToString() != null && record.getName().toString()
						.equals(hostname)){
					return record;
				}
			}
		}
		return null;
	}

	@Test
	public void whileSigningZones(){
		// This one verifies that DNS lookups fail while the zoneCache is being signed
		Runnable signer = new Runnable() {
			public void run(){
				final LoadingCache<ZoneKey, Zone> emptyZc =
						ZoneManager.createZoneCache(ZoneManager.ZoneCacheType.STATIC);
				Whitebox.setInternalState(ZoneManager.class, "zoneCache", emptyZc);
				try{
					Thread.sleep(200);
				} catch (InterruptedException ie) {
					// ehh
				}
				ZoneManager.initZoneCache(trafficRouter);
			}
		};

		Thread signThread = new Thread(signer);
		int resolveCount = 0;
		int loopCount = resolveCount;
		signThread.start();
		List<InetRecord> records = null;
		while (signThread.getState().compareTo(Thread.State.TERMINATED) != 0){
			loopCount++;
			records = zoneManager.resolve("edge-cache-001.dns-test.thecdn.example.com.");
			if (records == null) {
				resolveCount++;
			}
			try{
				Thread.sleep(100);
			} catch (InterruptedException ie) {
				// ehh
			}
		}
		assertThat("resolveCount should be greater than 0", resolveCount, greaterThan(0));
		assertThat("resolveCount should be less than loopCount", loopCount, greaterThanOrEqualTo(resolveCount));
		records = zoneManager.resolve("edge-cache-001.dns-test.thecdn.example.com.");
		assertThat("There should be 4 Records", records.size(), equalTo(4));
	}
}

