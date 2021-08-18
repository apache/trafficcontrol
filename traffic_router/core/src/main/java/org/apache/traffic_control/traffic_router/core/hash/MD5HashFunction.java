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

package org.apache.traffic_control.traffic_router.core.hash;

import java.math.BigInteger;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;

import org.springframework.stereotype.Component;

/**
 * For use with the Consistent Hash Algorithm using Java's
 * hashCode() method on a string value.
 */
@Component
public class MD5HashFunction {

    public double hash(final String value) {
        final byte[] valueBytes = value != null ? value.getBytes() : "".getBytes();
        return new BigInteger(1, md5Digest().digest(valueBytes)).doubleValue();
    }

    @SuppressWarnings("PMD.AvoidThrowingRawExceptionTypes")
    MessageDigest md5Digest() {
        // https://docs.oracle.com/javase/8/docs/api/java/security/MessageDigest.html

        // Every implementation of the Java platform is required to support the following standard MessageDigest algorithms:
        //
        // MD5
        // SHA-1
        // SHA-256

        try {
            return MessageDigest.getInstance("MD5");
        } catch (NoSuchAlgorithmException e) {
            // This should NEVER happen
            throw new RuntimeException("Failed to get MD5 message digest, something's very wrong!", e);
        }
    }
}
