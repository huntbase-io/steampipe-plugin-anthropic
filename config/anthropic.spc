connection "anthropic" {
  plugin = "local/anthropic"

  # Standard API key (sk-ant-api...), used for data-plane tables:
  # anthropic_model, anthropic_message_batch, anthropic_file.
  # Can also be set with the ANTHROPIC_API_KEY environment variable.
  # api_key = "sk-ant-api03-..."

  # Admin API key (sk-ant-admin...), used for organization tables:
  # anthropic_workspace, anthropic_workspace_member, anthropic_organization_member,
  # anthropic_invite, anthropic_api_key. Also accepted by anthropic_compliance_activity.
  # Can also be set with the ANTHROPIC_ADMIN_KEY environment variable.
  # admin_api_key = "sk-ant-admin01-..."

  # Compliance Access Key (sk-ant-api01..., created in claude.ai; Claude Enterprise only),
  # used for the anthropic_compliance_* tables. Scopes required per table:
  #   read:compliance_activities - anthropic_compliance_activity
  #   read:compliance_org_data   - organizations, roles, groups, settings
  #   read:compliance_user_data  - users, group members, chats, chat messages, projects, attachments
  # Can also be set with the ANTHROPIC_COMPLIANCE_ACCESS_KEY environment variable.
  # compliance_api_key = "sk-ant-api01-..."

  # Analytics API key (created in claude.ai > Organization settings > API by the
  # primary owner; read:analytics scope; Claude Enterprise only), used for the
  # anthropic_analytics_* tables. Note: anthropic_claude_code_analytics uses the
  # admin_api_key instead.
  # Can also be set with the ANTHROPIC_ANALYTICS_API_KEY environment variable.
  # analytics_api_key = "..."
}
