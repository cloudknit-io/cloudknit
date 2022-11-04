/*
 *  Copyright 2020 Curity AB
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

const util = require('util');
const { Strategy, TokenSet } = require('openid-client');
const base64url = require('base64url');

function CurityStrategy(options, verify) {
  const fallbackToUserInfoRequest = options.fallbackToUserInfoRequest ? options.fallbackToUserInfoRequest : false;
  // Always use PKCE with Curity Identity Server
  options.usePKCE = false;

  const decodeProfileFromIDToken = (accessToken, cb) => {
    const accessTokenValue = accessToken instanceof TokenSet && (accessToken.access_token || null) || accessToken;
    const refreshTokenValue = accessToken instanceof TokenSet && accessToken.refresh_token || null;
    const idTokenValue = accessToken instanceof TokenSet && accessToken.id_token || null;

    return Promise.resolve().then(() =>
      idTokenValue ? getIdTokenClaims(idTokenValue) : null
    ).then((profile) =>
      profile !== null ? Promise.resolve(profile) : (fallbackToUserInfoRequest? options.client.userinfo(accessToken) : Promise.resolve({}))
    ).then((profile) => {
      verify(accessTokenValue, refreshTokenValue, profile, cb);
    });
  };

  this._base = Object.getPrototypeOf(CurityStrategy.prototype);
  this._base.constructor.call(this, options, decodeProfileFromIDToken);

  this.name = 'curity';
}

const getIdTokenClaims = (idTokenValue) => JSON.parse(base64url.decode(idTokenValue.split('.')[1]));

util.inherits(CurityStrategy, Strategy);

module.exports = CurityStrategy;
