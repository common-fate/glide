import { ProviderSetup } from "src/utils/backend-client/types";

export const formatValidationErrorToText = (errMsgArray: ProviderSetup["configValidation"]) => {
  const result = errMsgArray.map((msg, index) => {

    const formattedLogs = msg.logs.map((log) =>
      `${log.level}: ${log.msg}`).join("\n");

    return `${index}:${msg.name}\n\t${formattedLogs}\n`
  })

  return result.join("")
}
