import React from "react";

export const UserLayout: React.FC<{ children?: React.ReactNode }> = ({
  children,
}) => {
  return <main>{children}</main>;
};
