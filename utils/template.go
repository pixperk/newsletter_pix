package utils

import "fmt"

// WrapInTemplate wraps the converted HTML body (including footer) in a
// dark, glassmorphic, Linear-inspired email template.
func WrapInTemplate(subject, htmlBody string) string {
	return fmt.Sprintf(`<!doctype html>
<html lang="en">
<head>
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
  <title>%s</title>
  <style media="all" type="text/css">
    @media only screen and (max-width: 640px) {
      .container { padding: 8px !important; width: 100%% !important; }
      .main { border-radius: 12px !important; }
      .wrapper { padding: 20px !important; }
      .wrapper p, .wrapper td, .wrapper span, .wrapper li { font-size: 15px !important; }
      pre { font-size: 13px !important; }
    }
    @media all {
      .ExternalClass { width: 100%%; }
      .ExternalClass, .ExternalClass p, .ExternalClass span, .ExternalClass font, .ExternalClass td, .ExternalClass div { line-height: 100%%; }
      .apple-link a { color: inherit !important; font-family: inherit !important; font-size: inherit !important; font-weight: inherit !important; line-height: inherit !important; text-decoration: none !important; }
      #MessageViewBody a { color: inherit; text-decoration: none; font-size: inherit; font-family: inherit; font-weight: inherit; line-height: inherit; }
    }

    /* Typography */
    .wrapper h1 { font-size: 26px; font-weight: 600; color: #ececec; margin: 0 0 20px 0; line-height: 1.3; letter-spacing: -0.5px; }
    .wrapper h2 { font-size: 21px; font-weight: 600; color: #dcdcdc; margin: 32px 0 12px 0; line-height: 1.3; letter-spacing: -0.3px; }
    .wrapper h3 { font-size: 17px; font-weight: 600; color: #cdcdcd; margin: 28px 0 10px 0; line-height: 1.4; }
    .wrapper p { font-size: 15px; font-weight: 400; color: #9a9a9a; margin: 0 0 16px 0; line-height: 1.75; }
    .wrapper a { color: #ccc; text-decoration: none; border-bottom: 1px solid rgba(255,255,255,0.15); transition: border-color 0.2s; }
    .wrapper ul, .wrapper ol { margin: 0 0 16px 0; padding-left: 20px; color: #9a9a9a; }
    .wrapper li { margin-bottom: 8px; line-height: 1.75; font-size: 15px; }
    .wrapper blockquote { margin: 20px 0; padding: 14px 20px; border-left: 2px solid rgba(255,255,255,0.15); background: rgba(255,255,255,0.02); color: #888; font-style: italic; border-radius: 0 6px 6px 0; }
    .wrapper hr { border: none; border-top: 1px solid rgba(255,255,255,0.06); margin: 32px 0; }
    .wrapper img { max-width: 100%%; height: auto; border-radius: 8px; }
    .wrapper strong { color: #d4d4d4; font-weight: 600; }
    .wrapper em { color: #888; }

    /* Code blocks */
    .wrapper pre { background: #0a0a0a; color: #c9d1d9; padding: 18px; border-radius: 10px; overflow-x: auto; font-size: 13.5px; line-height: 1.6; margin: 0 0 16px 0; border: 1px solid rgba(255,255,255,0.05); }
    .wrapper code { font-family: 'SF Mono', 'Fira Code', Consolas, 'Liberation Mono', Menlo, monospace; font-size: 13.5px; }
    .wrapper :not(pre) > code { background: rgba(255,255,255,0.06); color: #bbb; padding: 2px 7px; border-radius: 4px; font-size: 13.5px; border: 1px solid rgba(255,255,255,0.08); }
  </style>
</head>
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Inter, Roboto, Helvetica, Arial, sans-serif; -webkit-font-smoothing: antialiased; font-size: 15px; line-height: 1.75; -ms-text-size-adjust: 100%%; -webkit-text-size-adjust: 100%%; background-color: #0a0a0a; margin: 0; padding: 0;">
  <table role="presentation" border="0" cellpadding="0" cellspacing="0" class="body" style="border-collapse: separate; background-color: #0a0a0a; width: 100%%;" width="100%%">
    <tr>
      <td style="font-size: 15px; vertical-align: top;" valign="top">&nbsp;</td>
      <td class="container" style="vertical-align: top; max-width: 600px; padding: 0; padding-top: 40px; padding-bottom: 40px; width: 600px; margin: 0 auto;" width="600" valign="top">
        <div class="content" style="box-sizing: border-box; display: block; margin: 0 auto; max-width: 600px; padding: 0;">

          <!-- HEADER -->
          <table role="presentation" border="0" cellpadding="0" cellspacing="0" style="width: 100%%; margin-bottom: 24px;" width="100%%">
            <tr>
              <td style="padding: 0 8px; text-align: left;" align="left">
                <span style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Inter, Roboto, sans-serif; font-size: 15px; font-weight: 600; color: #e0e0e0; letter-spacing: -0.3px;">pixperk</span>
                <span style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Inter, Roboto, sans-serif; font-size: 15px; font-weight: 300; color: #333;">·</span>
                <span style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Inter, Roboto, sans-serif; font-size: 13px; font-weight: 400; color: #484848;">newsletter</span>
              </td>
            </tr>
          </table>

          <!-- MAIN CARD -->
          <table role="presentation" border="0" cellpadding="0" cellspacing="0" class="main" style="border-collapse: separate; background-color: #111111; border: 1px solid rgba(255,255,255,0.06); border-radius: 14px; width: 100%%;" width="100%%">
            <tr>
              <td class="wrapper" style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Inter, Roboto, Helvetica, Arial, sans-serif; font-size: 15px; vertical-align: top; box-sizing: border-box; padding: 36px;" valign="top">
                %s
              </td>
            </tr>
          </table>

          <!-- BOTTOM -->
          <table role="presentation" border="0" cellpadding="0" cellspacing="0" style="width: 100%%; margin-top: 20px;" width="100%%">
            <tr>
              <td style="text-align: center; padding: 0 8px;" align="center">
                <a href="https://www.pixperk.tech" style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Inter, Roboto, sans-serif; font-size: 12px; color: #3a3a3a; text-decoration: none;">pixperk.tech</a>
              </td>
            </tr>
          </table>

        </div>
      </td>
      <td style="font-size: 15px; vertical-align: top;" valign="top">&nbsp;</td>
    </tr>
  </table>
</body>
</html>`, subject, htmlBody)
}
