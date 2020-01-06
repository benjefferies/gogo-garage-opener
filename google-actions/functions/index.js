const { smarthome } = require("actions-on-google");
const functions = require("firebase-functions");
const axios = require("axios").default;

const app = smarthome();

function getAccessToken(headers) {
  const authorization = headers.authorization;
  console.log(`Authorization header: ${authorization}`);
  return authorization.substr(7);
}

async function getGarageState(accessToken) {
  const response = await axios.get(
    "https://open.mygaragedoor.space/garage/state",
    { headers: { Authorization: `Bearer ${accessToken}` } }
  );
  console.log(response);
  return response.data["Description"];
}

async function getUserId(accessToken) {
  const response = await axios.get(
    "https://gogo-garage-opener.eu.auth0.com/userinfo",
    { headers: { Authorization: `Bearer ${accessToken}` } }
  );
  console.log(response);
  return response.data["email"];
}

async function toggleGarageDoor(accessToken) {
  const response = await axios.post(
    "https://open.mygaragedoor.space/garage/toggle",
    {},
    { headers: { Authorization: `Bearer ${accessToken}` } }
  );
  console.log(response);
}

app.onSync(async (body, headers) => {
  const accessToken = getAccessToken(headers);
  const userId = await getUserId(accessToken)
  return {
    requestId: body.requestId,
    payload: {
      agentUserId: userId,
      devices: [
        {
          id: "garage-opener",
          type: "action.devices.types.GARAGE",
          traits: ["action.devices.traits.OpenClose"],
          name: {
            defaultNames: ["Bens Garage door"],
            name: "Bens Garage door",
            nicknames: ["Bens Garage door"]
          },
          willReportState: false,
          deviceInfo: {
            manufacturer: "Echosoft",
            model: "rpi",
            hwVersion: "1.0",
            swVersion: "1.0"
          }
        }
      ]
    }
  };
});

app.onQuery(async (body, headers) => {
  const accessToken = getAccessToken(headers);
  const userId = await getUserId(accessToken)
  var online = true;
  try {
    state = await getGarageState(accessToken);
  } catch (error) {
    online = false;
  }
  var openPercent = 100.0;
  if (state === "Closed") {
    openPercent = 0.0;
  }
  return {
    requestId: body.requestId,
    agentUserId: userId,
    payload: {
      devices: {
        "garage-opener": {
          on: online,
          online: online,
          openState: [
            {
              openPercent: openPercent
            }
          ]
        }
      }
    }
  };
});

app.onExecute(async (body, headers) => {
  const accessToken = getAccessToken(headers);
  const userId = await getUserId(accessToken)
  var online = true;
  try {
    const state = await getGarageState(accessToken);
    var openPercent = 100.0;
    if (state === "Open") {
      openPercent = 0.0;
    }
    await toggleGarageDoor(accessToken);
  } catch (error) {
    online = false;
  }
  return {
    requestId: body.requestId,
    agentUserId: userId,
    payload: {
      commands: [
        {
          ids: ["garage-opener"],
          status: "SUCCESS",
          states: {
            openPercent: openPercent,
            online: online
          }
        }
      ]
    }
  };
});

exports.fulfillment = functions.https.onRequest(app);
