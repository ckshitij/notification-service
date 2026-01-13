# Notification Service

## Template Module

The Template module is responsible for managing notification templates in a safe, versioned, and extensible way.
It provides APIs to create templates, manage versions, and render templates for different notification channels such as Email, Slack, and In-App.

This module is designed to work independently today and can be easily extracted into a separate service in the future.

### Core Concepts
1. Template
    - A Template represents a notification definition for a specific channel.
    - Identity
        - `(name + channel)` like `onboard_user + email`, `onboard_user + slack`
    - Templates contain metadata only, not message content.
    - Key fields
        - `name`
        - `description`
        - `channel` **(email / slack / in_app)**
        - `type` **(system / user)**
        - `active_version`
        - `timestamps`

2. Template Version
    - A TemplateVersion represents the actual message content.
    - Characteristics
        - Immutable
        - Versioned (1, 2, 3, …)
        - Only one version is active at a time
        - Linked to a template via template_id
    - Key fields
        - `version`
        - `subject` (optional, email-only)
        - `body`
        - `is_active`
        - `timestamps`

3. Renderer

    - The Renderer converts a template version into final output using runtime data.
        - Uses **Go** `text/template`
        - Stateless and pure
        - Fails fast on missing variables
        - Produces a channel-agnostic output:
            - `subject`
            - `body`
        - The renderer has no database or network dependencies.

### Template Lifecycle Flow

1. Create Template `POST /templates`
    - Creates a user-defined template
    - System templates are created only via migrations
    - Template is created without content

2. Add Template Version `POST /templates/{channel}/{name}/versions`
    - Creates a new immutable version
    - Automatically increments version number
    - Deactivates the previous active version
    - Updates active_version on the template

3. List Versions (Audit & Debug) `GET /templates/{channel}/{name}/versions`
    - Returns all versions (active + inactive)
    - Used by admin/UI/debugging tools
    - Notifications do not call this API

4. Render Template (Preview / Representation) `POST /templates/{channel}/{name}`
    - Renders the active version
    - Accepts runtime data as input
    - Returns final rendered output
    - No side effects (does not send notifications)

5. List Templates with Active Version (Summary) `GET /templates/summary`
    - Uses SQL JOIN to fetch templates + active version
    - Intended for admin/UI/ops use
    - Returns a read model, not a write model