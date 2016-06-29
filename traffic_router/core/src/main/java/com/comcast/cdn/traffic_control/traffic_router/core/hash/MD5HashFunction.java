/*
 * Copyright 2015 Comcast Cable Communications Management, LLC
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

package com.comcast.cdn.traffic_control.traffic_router.core.hash;

import java.math.BigInteger;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;

import org.apache.log4j.Logger;
import org.springframework.stereotype.Component;

/**
 * For use with the Consistent Hash Algorithm using Java's
 * hashCode() method on a string value.
 */
@Component
public class MD5HashFunction {
    private static final Logger LOGGER = Logger.getLogger(MD5HashFunction.class);

    private MessageDigest md5;

    public MD5HashFunction() {
        try {
            md5 = MessageDigest.getInstance("MD5");
        } catch (final NoSuchAlgorithmException e) {
            LOGGER.error(e.getMessage(), e);
        }
    }

    public double hash(final String value) {
        final BigInteger bi = new BigInteger(1, md5.digest(value != null ? value.getBytes() : "".getBytes()));
        return bi.doubleValue();
    }

}
