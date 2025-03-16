FROM golang AS go

ARG goproxy="https://proxy.golang.org,direct"

WORKDIR /rodent
COPY go.mod go.sum ./

ENV GOPROXY=$goproxy
ENV GO111MODULE=on
ENV CGO_ENABLED=0

RUN go mod download

COPY . .

RUN go build -ldflags "-s -w" -o rodent .

FROM debian:bookworm-slim
LABEL name="rodent" \
	maintainer="Yewolf <contact@yewolf.fr>" \
	description="Rodent's container"

# Install deps + add Chrome Stable + purge all the things
RUN apt-get update && apt-get install -y \
	apt-transport-https \
	ca-certificates \
	curl \
	gnupg \
	--no-install-recommends \
	&& apt-get update && apt-get install -y \
    chromium \
    chromium-l10n \
	fontconfig \
	fonts-ipafont-gothic \
	fonts-wqy-zenhei \
	fonts-thai-tlwg \
	fonts-kacst \
	fonts-symbola \
	fonts-noto \
	fonts-freefont-ttf \
	--no-install-recommends \
	&& apt-get purge --auto-remove -y curl gnupg \
	&& rm -rf /var/lib/apt/lists/*

# Add Chrome as a user
RUN groupadd -r user && useradd -r -g user -G audio,video user \
	&& mkdir -p /home/user && chown -R user:user /home/user

ENV BROWSER_PATH=/usr/bin/chromium

# Run as non-privileged
USER user

COPY --from=go /rodent/rodent /usr/bin/rodent

ENTRYPOINT [ "rodent" ]
