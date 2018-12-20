package sso

import (
	"github.com/skygeario/skygear-server/pkg/core/template"
)

type AuthHandlerHTMLProvider struct {
	APIEndPoint string
	JSSDKCDNURL string
}

func NewAuthHandlerHTMLProvider(APIEndPoint string, JSSDKCDNURL string) AuthHandlerHTMLProvider {
	return AuthHandlerHTMLProvider{
		APIEndPoint: APIEndPoint,
		JSSDKCDNURL: JSSDKCDNURL,
	}
}

func (i *AuthHandlerHTMLProvider) HTML() (out string, err error) {
	const templateString = `
<!DOCTYPE html>
<html>
<head>
<meta name=viewport content="width=device-width, initial-scale=1">
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<script type="text/javascript" src="{{ js_sdk_cdn_url }}"></script>
<script type="text/javascript">
skygear.pubsub.autoPubsub = false;
skygear.config({
	'endPoint': '{{ api_endpoint }}',
	'apiKey': '-'
}).then(function(container) {
	skygear.auth.oauthHandler();
}, function(err) {
	console.error(err);
});
</script>
</head>
<body>
</body>
</html>
	`
	context := map[string]interface{}{
		"api_endpoint":   i.APIEndPoint,
		"js_sdk_cdn_url": i.JSSDKCDNURL,
	}

	return template.ParseHTMLTemplate(templateString, context)
}
