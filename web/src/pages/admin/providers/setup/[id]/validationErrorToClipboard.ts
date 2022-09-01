import { ProviderSetup } from "src/utils/backend-client/types";

export const formatValidationErrorToText = (errMsgArray: ProviderSetup["configValidation"]) => {
    const result =  errMsgArray.map((msg, index) => {

      let formattedLogs = msg.logs.map((log) => 
      `${log.level}: ${log.msg}`).join("\n");
      
      const requiredString = `${index}:${msg.name}\n\t${formattedLogs}\n`

      return requiredString
    })

    return result.join("")
  }
