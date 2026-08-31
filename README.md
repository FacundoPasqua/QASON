# QASON — QA Agentic Testing, edición educativa

QASON convierte tu Claude Code en un equipo de QA con **tres agentes
especializados** que trabajan en cadena. Es la edición educativa de
[QATES](https://github.com/FacundoPasqua/Qates): mismos prompts, mismos
skills, recortado a lo esencial para aprender **IA agéntica aplicada a
QA** sin ruido.

## Qué vas a aprender usándolo

1. **Orquestación**: un agente coordinador que NO hace el trabajo — decide
   quién lo hace y sintetiza los resultados.
2. **Especialización**: tres sub-agentes con roles acotados (analista,
   diseñador de casos, automatizador) que superan a un generalista.
3. **Pipeline**: el flujo Spec-to-Test — de un ticket ambiguo a tests
   ejecutables, pasando por análisis de requisitos y diseño de casos.
4. **Skills**: capacidades atómicas en markdown que los agentes cargan
   según la tarea. Abrí cualquier `SKILL.md` y leelo: el prompt ES el
   código.

## Requisitos

- [Claude Code](https://claude.com/claude-code) instalado y autenticado
- Go 1.26+ (solo para compilar el instalador) — **instalalo antes de la clase**

## Instalación

```bash
git clone https://github.com/FacundoPasqua/Qason.git
cd Qason
go run ./cmd/qason        # wizard interactivo (con medialuna 🥐)
# o directo, sin wizard:
go run ./cmd/qason install
```

Eso instala en tu `~/.claude`:

| Qué | Dónde |
|---|---|
| 3 sub-agentes QA | `agents/qa-{analyst,test-designer,automator}.md` |
| 31 skills | `skills/qason/<skill>/SKILL.md` |
| Orquestador | bloque manejado dentro de `CLAUDE.md` |

Para desinstalar: `go run ./cmd/qason uninstall` (tu contenido propio de
`CLAUDE.md` queda intacto).

### ¿Usás GitHub Copilot en vez de Claude Code?

También funciona. VS Code ya busca sub-agentes en `~/.claude/agents` (lee el
formato de Claude) y skills en `~/.claude/skills`, así que el install de arriba
te deja los tres agentes y las 31 skills listos sin hacer nada más.

Lo único que Copilot **no** lee es `CLAUDE.md`, que es justo donde vive el
orquestador. Sin él tenés los especialistas pero nadie los encadena — y el
pipeline es la clase entera. Pasale la raíz de tu proyecto:

```bash
go run ./cmd/qason install --copilot .
```

Eso escribe el orquestador en `.github/copilot-instructions.md`, respetando las
reglas que el archivo ya tuviera. `uninstall` lo saca y deja las tuyas intactas.

Dos diferencias honestas: el campo `color` de los agentes no está entre los que
Copilot soporta (los colores del pipeline probablemente no se vean), y `model:
sonnet` usa nombres de modelo de Claude Code. Probalo antes de depender de ello.

## Primer ejercicio

Abrí Claude Code en cualquier proyecto y pegá un ticket — por ejemplo:

> Como usuario quiero recuperar mi contraseña por email.
> Criterios: un link "olvidé mi contraseña", se envía un mail con link
> de reseteo, y si el email no existe se muestra "No existe cuenta para
> este email".
>
> Analizá este ticket y generá los tests.

Mirá lo que pasa: el orquestador delega al `qa-analyst` (que va a
encontrar el problema de seguridad escondido en ese ticket — sí, hay
uno), pasa el análisis al `qa-test-designer`, y el `qa-automator`
escribe y CORRE los tests. Tres agentes, un pipeline, cero magia.

## Los tres agentes

| Agente | Color | Rol | Skills clave |
|---|---|---|---|
| `qa-analyst` | 🔵 cyan | Análisis de requisitos, test plans, matrices de riesgo | prd-analyzer, test-plan-gen, risk-matrix-gen |
| `qa-test-designer` | 🟢 green | Diseño de casos: funcionales, edge, negativos, exploratorios | e2e-test-gen, api-test-gen, exploratory-guide, test-data-gen |
| `qa-automator` | 🟣 purple | Código de automatización: unit, integration, e2e, performance, a11y | unit-test-gen, playwright-scaffold, output-validator |

Los colores no son decoración: cuando corre el pipeline, la task list de
Claude Code muestra tres agentes de colores distintos pasándose el trabajo.
La especialización se ve, no hay que explicarla.

## Cómo funciona por dentro

No hay servidor, no hay API key propia, no hay magia: QASON copia
archivos markdown a tu configuración de Claude Code. El instalador
completo son ~200 líneas de Go legibles ([installer.go](internal/installer/installer.go))
— leelo, es parte del material: vas a ver frontmatter YAML, un swap
atómico de directorios y un bloque manejado con marcadores. El binario
no hace ninguna llamada de red.

## Relación con QATES

QASON es un subconjunto de [QATES](https://github.com/FacundoPasqua/Qates),
el ecosistema completo (5 agentes, 51 skills, workflows, enforcement por
hooks, métricas de adopción y regresión de prompts). Lo que aprendés acá
transfiere 1:1. Los prompts se sincronizan desde QATES — los issues de
contenido van allá.

## Para alumnos

La guía completa — alcances, límites, los tres agentes en detalle y
cinco ejercicios prácticos — está en
[docs/GUIA-ALUMNOS.md](docs/GUIA-ALUMNOS.md). Empezá por ahí.

## Licencia

MIT — usalo, modificalo y enseñá con él.
