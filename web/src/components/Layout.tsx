import React from "react";
import { AdminNavbar } from "./nav/AdminNavbar";
import { Navbar } from "./nav/Navbar";

export const AdminLayout: React.FC<{ children?: React.ReactNode }> = ({
  children,
}) => {
  return (
    <main>
      <AdminNavbar />
      {children}
    </main>
  );
};

export const UserLayout: React.FC<{ children?: React.ReactNode }> = ({
  children,
}) => {
  return (
    <main>
      <Navbar />
      {children}
    </main>
  );
};
