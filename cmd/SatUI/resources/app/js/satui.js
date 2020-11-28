
async function sendMessage(functionName, functionParams) {
  const data = {
    name: functionName,
    payload: Array.from(functionParams),
  };
  console.log("Sending Data: ", data);
  return new Promise((resolve, reject) => {
    astilectron.sendMessage(data, function(message) {
      console.log("Received message", message, ...arguments);
      if (message) {
        if (message.name === 'error') {
          reject(message.payload);
        } else {
          resolve(message.payload);
        }
      } else {
        resolve();
      }
    });
  });
}

async function SatHelperApp_SaveConfig() {
  return sendMessage('SatHelperApp_SaveConfig', arguments);
}

async function SatHelperApp_IsConfigLoaded() {
  return sendMessage('SatHelperApp_IsConfigLoaded', arguments);
}

async function SatHelperApp_GetConfig() {
  return sendMessage('SatHelperApp_GetConfig', arguments);
}

async function SatHelperApp_SetConfig(config) {
  return sendMessage('SatHelperApp_SetConfig', arguments);
}

async function SatHelperApp_LoadConfig() {
  return sendMessage('SatHelperApp_LoadConfig', arguments);
}

async function SatHelperApp_StartServer() {
  return sendMessage('SatHelperApp_StartServer', arguments);
}

async function SatHelperApp_StopServer() {
  return sendMessage('SatHelperApp_StopServer', arguments);
}

async function SatHelperApp_Exit() {
  return sendMessage('SatHelperApp_Exit', arguments);
}

async function SatHelperApp_ServerIsRunning() {
  return sendMessage('SatHelperApp_ServerIsRunning', arguments);
}

const Terminal = require("./terminal/index.js");
const terminal = new Terminal({columns: 80, rows: 16});
