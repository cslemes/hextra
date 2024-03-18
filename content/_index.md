---
title: My Site
toc: false
---

This is the landing page.

## Explore

{{< cards >}}
  {{< card link="docs" title="Docs" icon="book-open" >}}
  {{< card link="about" title="About" icon="user" >}}
{{< /cards >}}


{{ range first 10 .Pages }}
  <article>
    <!-- this <div> includes the title summary -->
    <div>
      <h2><a href="{{ .RelPermalink }}">{{ .Title }}</a></h2>
      {{ .Summary }}
    </div>
    {{ if .Truncated }}
      <!-- This <div> includes a read more link, but only if the summary is truncated... -->
      <div>
        <a href="{{ .RelPermalink }}">Read More…</a>
      </div>
    {{ end }}
  </article>
{{ end }}


## Documentation

For more information, visit [Hextra](https://imfing.github.io/hextra).
