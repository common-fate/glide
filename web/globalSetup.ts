import { FullConfig } from "@playwright/test";
import dotenv from "dotenv";

// Read from default ".env" file.
dotenv.config();

async function globalSetup(config: FullConfig) {}
export default globalSetup;
