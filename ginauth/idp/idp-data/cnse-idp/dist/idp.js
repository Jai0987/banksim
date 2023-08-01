import { OAuth2Server } from 'oauth2-mock-server';
import * as fs from 'fs';
//Get Environment Variables and Set Defaults
const PORT = process.env.IDP_PORT || 1234;
const HOST = process.env.IDP_HOST || '0.0.0.0';
const CONFIG = process.env.IDP_CONFIG || './idp-cfg.json';
const AUDIENCE = process.env.IDP_AUDIENCE || 'cci.drexel.edu/api';
const DEFAULT_CLIENT = {
    client_id: 'guest-client',
    client_secret: undefined,
    account: {
        id: 0,
        group: '*',
        role: 'guest'
    }
};
const DEFAULT_USER = {
    subject: 'guest-user',
    account: {
        id: 0,
        group: '*',
        role: 'guest'
    }
};
console.log("Config:", HOST, PORT, CONFIG);
//Read and Process the config file
const jsonString = fs.readFileSync(CONFIG, 'utf-8');
const config = JSON.parse(jsonString);
let server = new OAuth2Server();
// Generate a new RSA key and add it to the keystore
await server.issuer.keys.generate('RS256');
// Start the server
await server.start(+PORT, HOST);
console.log('Issuer URL:', server.issuer.url); // -> http://localhost:8080
//Register a callback to alter the token prior to it being signed
server.service.on('beforeTokenSigning', (token, req) => {
    //UNCOMMENT TO DO SOME DEBUGGING
    //console.log('beforeTokenSigning/t', token);
    //console.log('beforeTokenSigning/r', req.body.grant_type, req.body.client_id);
    //add the audience - aud -  to the token
    token.payload.aud = AUDIENCE;
    switch (req.body.grant_type) {
        case 'client_credentials':
            buildClientToken(token, req);
            break;
        case 'password':
            buildUserToken(token, req);
            break;
        case 'authorization_code':
            buildUserToken(token, req);
            break;
        default:
            throw new Error('Unsupported grant type');
    }
});
function buildClientToken(token, req) {
    const clientId = req.body.client_id;
    if (clientId) {
        const client = config.find((idpConfigItem) => idpConfigItem.client_id === clientId);
        if (client) {
            token.payload.cci = sanitizeClient(client);
        }
        else {
            token.payload.cci = DEFAULT_CLIENT;
        }
    }
}
function buildUserToken(token, req) {
    const username = req.body.username;
    if (username) {
        //modify sub
        token.payload.sub = username;
        //modify scopes make more realistic
        token.payload.scope = ["openid", "profile", "email", "offline_access"];
        const user = config.find((idpConfigItem) => idpConfigItem.subject === username);
        if (user) {
            token.payload.cci = user;
        }
        else {
            token.payload.cci = DEFAULT_USER;
        }
    }
}
//This function makes sure the client secret is not leaked out
function sanitizeClient(client) {
    return { client_id: client.client_id, account: client.account };
}
//# sourceMappingURL=idp.js.map