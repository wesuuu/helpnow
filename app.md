I want to create an application that allows business people to manage AIs for specific departments, meaning all the departments that are required to do business but are not necessarily core to the business. Start first with a marketing agent that we can give tools and routines to.

Marketing agent should be able to:
- Create email campaigns
- Create social media posts
- Create blogs
- Create ads
- Give advice on SEO

Built-in tools:
- Email marketing
- Social media marketing
- Blog marketing
- Ad marketing
- SEO marketing

Users should be able to create routines based off of the built-in tools.

Make sure the architecture is open so we can add more agents and they can share tools and routines. Routines should also allow the user to setup human-in-the-loop workflows, so that we can checkpoint the AI's work and have a human review it before it goes live.

Users should also be allowed to select which AI model they want to use, with a list of built-ins, but also a bring your own key option so they can use whatever they already have configured.

I want the main code to be in Python since it has the best AI libraries, use Golang for the backend API with echo, the UI should be in svelte. The python AI code should communicate with the API via gRPC.

The data model should allow for access controls:
- Organization -> Teams -> Users
- Roles: read and write access depending on team/user

The app should also have a report dashboard that shows statistics about what the AI has done, including an event log.

The app should be containerized also.

Use SQLalchemy and postgres for the database and design the data model to include all the above requirements.

When writing API routes for the go backend, use the built-in SQL driver.

Use strands from aws for the python AI agent module

Write unit tests where it makes sense. After the app is working, we'll work on creating integration tests for the rest.


##

Write me the frontend svelte code for the app. Use tailwindcss for styling.

Make it so we can have a light and dark mode.

Implement the following pages:
- login
- signup
- register
- dashboard
- organization settings
- team settings
- user settings

Once logged in, have the same shared layout with a sidebar to the right and a nav header at the top.




Create a new directory called integrations and write me an npm package under clients that can collect marketing statistics like impressions, click rate, any metric that is useful for sales page. Make it so we can automatically calculate things like impressions and have methods or functions to record when they do things like convert or click. All records should be sent to the backend of helpnow, and helpnow should get the organization_id from a db lookup 


## Site

Analytics and workflow builder for marketing. They need little involvement from IT other than access controls to specific databases, apps, etc.


Workflows include:
- Lead generation
- Email marketing
- Social media marketing
  - Influencer research
  - Influencer outreach
  - Social media posting
- Blog marketing
- Ad marketing
- SEO marketing
- Enrichment

Analytics include:
- Impressions
- Clicks
- Conversions
- Cost

Tightly integrated with platforms like instagram, facebook, twitter, shopify, etc.

Workflows can deliver customer data to other apps and databases. 

## Agents

- researcher: find influencer contact info, send intro emails, schedule followups
- content creation:



## Standard workflow actions

- Send Ad
- Send Email
- Apply Coupon
- Send Feedback
- Send GDPR Request
- Send Webhook
- Send Slack Message
- Send SMS


We need to update the agent objects to have a configuration section. This will be a json which can store things required for model providers and since this configuration may contain sensitive information like API_KEYS, you should store them in vault 


## Specific workflows

- Campaigns
- Audience Segmentation & Enrichment

## Plugins

Categorize by:
- Enrichment
- Data Sources
- eCommerce
- Social Media
- Email
- SEO
- CRM

## Tabs

- Campaigns
  - Email
  - Ad
  - Social Media
  - Blog
  - SEO

- Audience Segmentation & Enrichment
  - Audiences
  - 
- Plugins

## Analytics

Attention metrics
- subscriptions
- ratings
- reach
- impressions
- engagement

Conversion metrics
- sales
- leads
- signups
- signins
- downloads

## Integrating various workflow actions

- Users submit requests for a workflow action
- AI bot codes this in the background
- Submits it to beta
- Users at their site can re-sync to beta to get the latest version of the workflow action

- Once beta works and users verify this feature is ok, it graduates to pre-production review
- after pre-production review, it graduates to production and it is just a part of the standard workflow actions