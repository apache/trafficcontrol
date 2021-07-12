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

package org.apache.traffic_control.traffic_router.core.util;

import com.fasterxml.jackson.databind.JsonNode;

public class JsonUtils {

    public static long getLong(final JsonNode jsonNode, final String key) throws JsonUtilsException {
        if (jsonNode == null || !jsonNode.has(key)) {
            throwException(key);
        }

        return jsonNode.get(key).asLong();
    }

    public static long optLong(final JsonNode jsonNode, final String key, final long d) {
        if (jsonNode == null || !jsonNode.has(key)) {
            return d;
        }

        return jsonNode.get(key).asLong(d);
    }

    public static long optLong(final JsonNode jsonNode, final String key) {
        return optLong(jsonNode, key, 0);
    }

    public static double getDouble(final JsonNode jsonNode, final String key) throws JsonUtilsException {
        if (jsonNode == null || !jsonNode.has(key)) {
            throwException(key);
        }

        return jsonNode.get(key).asDouble();
    }

    public static double optDouble(final JsonNode jsonNode, final String key, final double d) {
        if (jsonNode == null || !jsonNode.has(key)) {
            return d;
        }

        return jsonNode.get(key).asDouble(d);
    }

    public static double optDouble(final JsonNode jsonNode, final String key) {
        return optDouble(jsonNode, key, 0);
    }

    public static int getInt(final JsonNode jsonNode, final String key) throws JsonUtilsException {
        if (jsonNode == null || !jsonNode.has(key)) {
            throwException(key);
        }

        return jsonNode.get(key).asInt();
    }

    public static int optInt(final JsonNode jsonNode, final String key, final int d) {
        if (jsonNode == null || !jsonNode.has(key)) {
            return d;
        }

        return jsonNode.get(key).asInt(d);
    }

    public static int optInt(final JsonNode jsonNode, final String key) {
        return optInt(jsonNode, key, 0);
    }

    public static boolean getBoolean(final JsonNode jsonNode, final String key) throws JsonUtilsException {
        if (jsonNode == null || !jsonNode.has(key)) {
            throwException(key);
        }

        return jsonNode.get(key).asBoolean();
    }

    public static boolean optBoolean(final JsonNode jsonNode, final String key, final boolean d) {
        if (jsonNode == null || !jsonNode.has(key)) {
            return d;
        }

        return jsonNode.get(key).asBoolean(d);
    }

    public static boolean optBoolean(final JsonNode jsonNode, final String key) {
        return optBoolean(jsonNode, key, false);
    }

    public static String getString(final JsonNode jsonNode, final String key) throws JsonUtilsException {
        if (jsonNode == null || !jsonNode.has(key)) {
            throwException(key);
        }

        return jsonNode.get(key).asText();
    }

    public static String optString(final JsonNode jsonNode, final String key, final String d) {
        if (jsonNode == null || !jsonNode.has(key)) {
            return d;
        }

        return jsonNode.get(key).asText(d);
    }

    public static String optString(final JsonNode jsonNode, final String key) {
        return optString(jsonNode, key, "");
    }

    public static JsonNode getJsonNode(final JsonNode jsonNode, final String key) throws JsonUtilsException {
        if (jsonNode == null || !jsonNode.has(key)) {
            throwException(key);
        }

        return jsonNode.get(key);
    }

    public static void throwException(final String key) throws JsonUtilsException {
        throw new JsonUtilsException("Failed querying JSON for key: " + key);
    }
}
