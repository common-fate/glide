import * as React from "react";
import { Days, DurationInput, Hours, Minutes, Weeks } from "./DurationInput";
import { test, expect } from "@playwright/experimental-ct-react";

test("minutes renders", async ({ mount }) => {
  const component = await mount(
    <DurationInput onChange={(v) => {}}>
      <Minutes />
    </DurationInput>
  );
  await expect(component).toContainText("mins");
});

test("hours renders", async ({ mount }) => {
  const component = await mount(
    <DurationInput onChange={(v) => {}}>
      <Hours />
    </DurationInput>
  );
  await expect(component).toContainText("hrs");
});
test("hours and minutes renders", async ({ mount }) => {
  const component = await mount(
    <DurationInput onChange={(v) => {}}>
      <Minutes />
      <Hours />
      <Days />
      <Weeks />
    </DurationInput>
  );
  await expect(component).toContainText("mins");
  await expect(component).toContainText("hrs");
  await expect(component).toContainText("days");
  await expect(component).toContainText("weeks");
});

test("minutes increment decrement", async ({ mount }) => {
  const component = await mount(
    <DurationInput onChange={(v) => {}}>
      <Minutes />
    </DurationInput>
  );

  await expect(component.locator("#minute-duration-input")).toHaveValue("0");
  await expect(component.locator("#decrement")).toBeDisabled();
  await component.locator(`#increment >> visible=true`).first().click();
  await expect(component.locator("#minute-duration-input")).toHaveValue("1");
  await expect(component.locator("#decrement")).toBeEnabled();
  await component.locator(`#decrement >> visible=true`).first().click();
  await expect(component.locator("#minute-duration-input")).toHaveValue("0");
  await expect(component.locator("#decrement")).toBeDisabled();
});

test("hours increment decrement", async ({ mount }) => {
  const component = await mount(
    <DurationInput onChange={(v) => {}}>
      <Hours />
    </DurationInput>
  );
  await expect(component.locator("#hour-duration-input")).toHaveValue("0");
  await expect(component.locator("#decrement")).toBeDisabled();
  await component.locator(`#increment >> visible=true`).first().click();
  await expect(component.locator("#hour-duration-input")).toHaveValue("1");
  await expect(component.locator("#decrement")).toBeEnabled();
  await component.locator(`#decrement >> visible=true`).first().click();
  await expect(component.locator("#hour-duration-input")).toHaveValue("0");
  await expect(component.locator("#decrement")).toBeDisabled();
});

test("days increment decrement", async ({ mount }) => {
  const component = await mount(
    <DurationInput onChange={(v) => {}}>
      <Days />
    </DurationInput>
  );
  await expect(component.locator("#day-duration-input")).toHaveValue("0");
  await expect(component.locator("#decrement")).toBeDisabled();
  await component.locator(`#increment >> visible=true`).first().click();
  await expect(component.locator("#day-duration-input")).toHaveValue("1");
  await expect(component.locator("#decrement")).toBeEnabled();
  await component.locator(`#decrement >> visible=true`).first().click();
  await expect(component.locator("#day-duration-input")).toHaveValue("0");
  await expect(component.locator("#decrement")).toBeDisabled();
});

test("weeks increment decrement", async ({ mount }) => {
  const component = await mount(
    <DurationInput onChange={(v) => {}}>
      <Weeks />
    </DurationInput>
  );
  await expect(component.locator("#week-duration-input")).toHaveValue("0");
  await expect(component.locator("#decrement")).toBeDisabled();
  await component.locator(`#increment >> visible=true`).first().click();
  await expect(component.locator("#week-duration-input")).toHaveValue("1");
  await expect(component.locator("#decrement")).toBeEnabled();
  await component.locator(`#decrement >> visible=true`).first().click();
  await expect(component.locator("#week-duration-input")).toHaveValue("0");
  await expect(component.locator("#decrement")).toBeDisabled();
});
