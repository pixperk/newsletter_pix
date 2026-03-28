package utils

import "fmt"

// WrapInTemplate wraps the converted HTML body (including footer) in a
// responsive, email-client-compatible HTML email template.
func WrapInTemplate(subject, htmlBody string) string {
	return fmt.Sprintf(`<!doctype html>
<html lang="en">
<head>
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
  <title>%s</title>
  <style media="all" type="text/css">
    @media only screen and (max-width: 640px) {
      .container { padding: 0 !important; width: 100%% !important; }
      .main { border-left-width: 0 !important; border-radius: 0 !important; border-right-width: 0 !important; }
      .wrapper { padding: 16px !important; }
      .main p, .main td, .main span, .main li { font-size: 16px !important; }
      pre { font-size: 13px !important; }
    }
    @media all {
      .ExternalClass { width: 100%%; }
      .ExternalClass, .ExternalClass p, .ExternalClass span, .ExternalClass font, .ExternalClass td, .ExternalClass div { line-height: 100%%; }
      .apple-link a { color: inherit !important; font-family: inherit !important; font-size: inherit !important; font-weight: inherit !important; line-height: inherit !important; text-decoration: none !important; }
      #MessageViewBody a { color: inherit; text-decoration: none; font-size: inherit; font-family: inherit; font-weight: inherit; line-height: inherit; }
    }

    /* Content typography */
    .wrapper h1 { font-size: 26px; font-weight: 700; color: #1a1a1a; margin: 0 0 16px 0; line-height: 1.3; }
    .wrapper h2 { font-size: 22px; font-weight: 700; color: #1a1a1a; margin: 28px 0 12px 0; line-height: 1.3; }
    .wrapper h3 { font-size: 18px; font-weight: 600; color: #1a1a1a; margin: 24px 0 10px 0; line-height: 1.3; }
    .wrapper p { font-size: 16px; font-weight: normal; color: #374151; margin: 0 0 16px 0; line-height: 1.7; }
    .wrapper a { color: #2563eb; text-decoration: underline; }
    .wrapper a:hover { color: #1d4ed8; }
    .wrapper ul, .wrapper ol { margin: 0 0 16px 0; padding-left: 24px; color: #374151; }
    .wrapper li { margin-bottom: 6px; line-height: 1.7; font-size: 16px; }
    .wrapper blockquote { margin: 16px 0; padding: 12px 20px; border-left: 4px solid #e5e7eb; background: #f9fafb; color: #4b5563; font-style: italic; }
    .wrapper hr { border: none; border-top: 1px solid #e5e7eb; margin: 28px 0; }
    .wrapper img { max-width: 100%%; height: auto; }

    /* Code blocks — preserve monokai from goldmark-highlighting */
    .wrapper pre { background: #272822; color: #f8f8f2; padding: 16px; border-radius: 8px; overflow-x: auto; font-size: 14px; line-height: 1.5; margin: 0 0 16px 0; }
    .wrapper code { font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace; font-size: 14px; }
    .wrapper :not(pre) > code { background: #f3f4f6; color: #e11d48; padding: 2px 6px; border-radius: 4px; font-size: 14px; }
    .wrapper strong { color: #111827; }
    .wrapper em { color: #4b5563; }
  </style>
</head>
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; -webkit-font-smoothing: antialiased; font-size: 16px; line-height: 1.7; -ms-text-size-adjust: 100%%; -webkit-text-size-adjust: 100%%; background-color: #f4f5f6; margin: 0; padding: 0;">
  <table role="presentation" border="0" cellpadding="0" cellspacing="0" class="body" style="border-collapse: separate; background-color: #f4f5f6; width: 100%%;" width="100%%">
    <tr>
      <td style="font-size: 16px; vertical-align: top;" valign="top">&nbsp;</td>
      <td class="container" style="vertical-align: top; max-width: 600px; padding: 0; padding-top: 32px; padding-bottom: 32px; width: 600px; margin: 0 auto;" width="600" valign="top">
        <div class="content" style="box-sizing: border-box; display: block; margin: 0 auto; max-width: 600px; padding: 0;">

          <!-- HEADER -->
          <table role="presentation" border="0" cellpadding="0" cellspacing="0" style="width: 100%%; margin-bottom: 16px;" width="100%%">
            <tr>
              <td style="padding: 0 4px; text-align: left;" align="left">
                <span style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; font-size: 20px; font-weight: 700; color: #111827; letter-spacing: -0.5px;">pixperk</span>
              </td>
              <td style="padding: 0 4px; text-align: right;" align="right">
                <span style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; font-size: 13px; color: #9ca3af;">newsletter</span>
              </td>
            </tr>
          </table>

          <!-- MAIN CARD -->
          <table role="presentation" border="0" cellpadding="0" cellspacing="0" class="main" style="border-collapse: separate; background: #ffffff; border: 1px solid #e5e7eb; border-radius: 12px; width: 100%%;" width="100%%">
            <tr>
              <td class="wrapper" style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; font-size: 16px; vertical-align: top; box-sizing: border-box; padding: 32px;" valign="top">
                %s
              </td>
            </tr>
          </table>

        </div>
      </td>
      <td style="font-size: 16px; vertical-align: top;" valign="top">&nbsp;</td>
    </tr>
  </table>
</body>
</html>`, subject, htmlBody)
}
