import { User } from "./backend-client/types";

export const userName = (user: User) => {
  if (user.firstName === "" && user.lastName === "") return "";
  return user.firstName + " " + user.lastName;
};
