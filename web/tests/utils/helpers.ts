import { Page } from "@playwright/test";

//helper functions to click elements that are visible 

export const fillFormElement = async (
    inputType: 'input' | 'textarea',
    name: string,
    value: string,
    page: Page
  )  => {
    
       await page
      .locator(`${inputType}[name=${name}] >> visible=true`)
      .fill(value);
  }

  export const clickFormElement = async (
    inputType: 'input' | 'textarea',
    name: string,
    page: Page
  )  => {
    
       await page
      .locator(`input:has-text("${name}") >> visible=true`)
      .click();
    
    
  }