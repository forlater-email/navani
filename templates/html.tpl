<!DOCTYPE HTML PUBLIC "-//W3C//DTD XHTML 1.0 Transitional //EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:v="urn:schemas-microsoft-com:vml" xmlns:o="urn:schemas-microsoft-com:office:office">
<head>
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta name="x-apple-disable-message-reformatting">
  
  <style type="text/css">
    img, video {
       align: center;
       object-fit: cover;
       width: 400px;
       height: 200px;
    }

    body {
      margin: 0;
      padding: 0;
      max-width: 640px;
    }
    a[x-apple-data-detectors='true'] {
      color: inherit !important;
      text-decoration: none !important;
    }

    .page {
        font-family: serif;
    }
  </style>
</head>

<body>
  <em>original link</em>: <a href="{{ .URL }}">{{ .URL }}</a>
  <h1>
    {{ .Title }}
  </h1>
  <h2>
    {{ .Byline }}
  </h2>
  <hr>
  {{ .Content }}
</body>

</html>
