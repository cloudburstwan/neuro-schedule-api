# [UNOFFICIAL] Neuro-sama Schedule API
This repository contains the source code for the **unofficial** Neuro-sama stream schedule API available at https://schedule-api.nwero.net/.

## How to use?
This API currently provides 2 data formats: JSON and XML (RSS). More formats can (probably) be added upon request.

To receive JSON via REST, make a GET request to `https://schedule-api.nwero.net/schedule` with an `Accept` header set to one of the following values:
- ` ` (Nothing, the API will default to JSON)
- `*/*` (Anything)
- `application/*` (Application data)
- `application/json` (JSON payloads)

> [!IMPORTANT]
> There are plans to implement a EventStream JSON API which will be available by setting `Accept` to `text/event-stream` (or `text/*`). As this is currently not implemented, attempts to request this content will return a 501 Not Implemented HTTP error code.

To receive XML (RSS), make a GET request to `https://schedule-api.nwero.net/schedule.xml`. There is no need to provide any additional headers for this request.

If a data format you require is not provided by this API, please take a look at [our instructions for contributing](#Contributing).

## How does it work?
This API is technically a fancy Discord wrapper, listening to a channel that follows the #schedule channel in Neuro-sama Headquarters.

We initially attempt to parse content with Regex (for schedule format parsing) and basic text matching (for streamers in common streams), however if parsing fails, it is handed off to an LLM via OpenAI's API.

> [!CAUTION]
> While testing has shown that the LLM fallback works well in most cases, please note that during the initial few weeks/months of this API's operation, incorrect data could be passed by the LLM to the output. If this occurs, create an [Issue](https://github.com/cloudburstwan/neuro-schedule-api/issues) reporting the incorrect data.

This means that while the API is unofficial, the data for it comes from an official source. Barring LLM-related issues, the data provided by this API endpoint can be trusted.

## Contributing
You are free to contribute to this project by submitting [Pull Requests](https://github.com/cloudburstwan/neuro-schedule-api/pulls), suggesting changes via [Issues](https://github.com/cloudburstwan/neuro-schedule-api/issues), or forking the project for your own purposes, as long as you follow our [License](https://github.com/cloudburstwan/neuro-schedule-api/blob/main/LICENSE).