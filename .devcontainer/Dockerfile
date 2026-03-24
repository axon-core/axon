FROM mcr.microsoft.com/devcontainers/go:1.25-bookworm

USER root

# kind (for local K8s clusters)
RUN curl -Lo /usr/local/bin/kind https://kind.sigs.k8s.io/dl/latest/kind-linux-amd64 \
  && chmod +x /usr/local/bin/kind

# Claude Code
RUN curl -fsSL https://claude.ai/install.sh | bash

# claude-bedrock helper (runs Claude Code via AWS Bedrock)
RUN <<'EOF'
cat >> /etc/zshrc << 'SHELL_FUNC'

# Claude / Bedrock helper function
claude-bedrock() {
  CLAUDE_CODE_USE_BEDROCK=1 AWS_REGION=us-west-2 claude "$@"
}
SHELL_FUNC
cat >> /etc/bash.bashrc << 'SHELL_FUNC'

# Claude / Bedrock helper function
claude-bedrock() {
  CLAUDE_CODE_USE_BEDROCK=1 AWS_REGION=us-west-2 claude "$@"
}
SHELL_FUNC
EOF

# Oh My Zsh
RUN sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)" "" --unattended

USER vscode
