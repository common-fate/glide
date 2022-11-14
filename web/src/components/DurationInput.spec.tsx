import { expect, test } from "@playwright/experimental-ct-react";
import * as React from "react";
import { Days, DurationInput, Hours, Minutes, Weeks } from "./DurationInput";

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

const MINUTE = 60;
const HOUR = 3600;
const DAY = 86400;
const WEEK = 7 * DAY;
const MONTH = 30 * DAY;

const time = 3 * WEEK + 3 * DAY + 2 * HOUR + 1 * MINUTE; // 2080860

test("duration with max of 3weeks, 3days, 2hours, 1min works as expected", async ({
  mount,
}) => {
  const component = await mount(
    <DurationInput
      min={60}
      max={time}
      defaultValue={time}
      value={time}
      onChange={(v) => {}}
    >
      <Weeks />
      <Days />
      <Hours />
      <Minutes />
    </DurationInput>
  );
  await expect(component.locator("#week-duration-input")).toHaveValue("3");
  await expect(component.locator("#day-duration-input")).toHaveValue("3");
  await expect(component.locator("#hour-duration-input")).toHaveValue("2");
  await expect(component.locator("#minute-duration-input")).toHaveValue("1");
});

test("actions work as expected", async ({ mount }) => {
  const component = await mount(
    <DurationInput
      min={60}
      max={time}
      defaultValue={time}
      value={time}
      onChange={(v) => {}}
    >
      <Weeks />
      <Days />
      <Hours />
      <Minutes />
    </DurationInput>
  );
  // DECREMENT
  // decrement week
  await component.locator(`#decrement >> visible=true`).first().click();
  await expect(component.locator("#week-duration-input")).toHaveValue("2");
  await expect(component.locator("#day-duration-input")).toHaveValue("3");
  await expect(component.locator("#hour-duration-input")).toHaveValue("2");
  await expect(component.locator("#minute-duration-input")).toHaveValue("1");
  // decrement day
  await component.locator(`#decrement >> visible=true`).nth(1).click();
  await expect(component.locator("#week-duration-input")).toHaveValue("2");
  await expect(component.locator("#day-duration-input")).toHaveValue("2");
  await expect(component.locator("#hour-duration-input")).toHaveValue("2");
  await expect(component.locator("#minute-duration-input")).toHaveValue("1");
  // decrement hour
  await component.locator(`#decrement >> visible=true`).nth(2).click();
  await expect(component.locator("#week-duration-input")).toHaveValue("2");
  await expect(component.locator("#day-duration-input")).toHaveValue("2");
  await expect(component.locator("#hour-duration-input")).toHaveValue("1");
  await expect(component.locator("#minute-duration-input")).toHaveValue("1");
  // decrement minute
  await component.locator(`#decrement >> visible=true`).nth(3).click();
  await expect(component.locator("#week-duration-input")).toHaveValue("2");
  await expect(component.locator("#day-duration-input")).toHaveValue("2");
  await expect(component.locator("#hour-duration-input")).toHaveValue("1");
  await expect(component.locator("#minute-duration-input")).toHaveValue("0");
  // INCREMENT
  // increment week
  await component.locator(`#increment >> visible=true`).first().click();
  await expect(component.locator("#week-duration-input")).toHaveValue("3");
  await expect(component.locator("#day-duration-input")).toHaveValue("2");
  await expect(component.locator("#hour-duration-input")).toHaveValue("1");
  await expect(component.locator("#minute-duration-input")).toHaveValue("0");
  // increment day
  await component.locator(`#increment >> visible=true`).nth(1).click();
  await expect(component.locator("#week-duration-input")).toHaveValue("3");
  await expect(component.locator("#day-duration-input")).toHaveValue("3");
  await expect(component.locator("#hour-duration-input")).toHaveValue("1");
  await expect(component.locator("#minute-duration-input")).toHaveValue("0");
  // increment hour
  await component.locator(`#increment >> visible=true`).nth(2).click();
  await expect(component.locator("#week-duration-input")).toHaveValue("3");
  await expect(component.locator("#day-duration-input")).toHaveValue("3");
  await expect(component.locator("#hour-duration-input")).toHaveValue("2");
  await expect(component.locator("#minute-duration-input")).toHaveValue("0");
  // increment minute
  await component.locator(`#increment >> visible=true`).nth(3).click();
  await expect(component.locator("#week-duration-input")).toHaveValue("3");
  await expect(component.locator("#day-duration-input")).toHaveValue("3");
  await expect(component.locator("#hour-duration-input")).toHaveValue("2");
  await expect(component.locator("#minute-duration-input")).toHaveValue("1");
});
