# QASON — Guía para alumnos

> **QA Agéntico, edición educativa.** Este documento es tu mapa: qué es
> QASON, qué podés hacer con él, qué NO hace, y los ejercicios para
> aprender a dirigir agentes de IA en QA.

---

## 1. Qué es QASON (y qué no es)

QASON convierte tu Claude Code en un **equipo de QA de tres agentes
especializados** coordinados por un orquestador. No es un chatbot que
"sabe de testing": es una arquitectura donde cada agente tiene un rol
acotado, un prompt propio y un set de skills, y el trabajo fluye entre
ellos como en un equipo real.

**La idea central que vas a aprender**: un agente especializado con
contexto acotado supera a un generalista que intenta hacer todo. Igual
que en un equipo humano.

### Alcances — qué SÍ hace

| Capacidad | Agente responsable |
|---|---|
| Analizar tickets/PRDs, detectar requisitos faltantes y ambigüedades | `qa-analyst` |
| Generar test plans y matrices de riesgo | `qa-analyst` |
| Diseñar casos de prueba (funcionales, edge, negativos, exploratorios) | `qa-test-designer` |
| Diseñar casos de API a partir de specs (OpenAPI, descripciones) | `qa-test-designer` |
| Generar datos de prueba y casos data-driven | `qa-test-designer` |
| Escribir código de automatización (unit, integration, e2e) | `qa-automator` |
| Ejecutar los tests generados y reportar resultados | `qa-automator` |
| Scaffolding de frameworks (Playwright, Cypress, Appium, Postman) | `qa-automator` |
| Encadenar todo lo anterior: de ticket a tests corriendo | orquestador |

### Límites — qué NO hace (leé esto dos veces)

- **No reemplaza tu criterio.** Los agentes producen borradores de
  calidad profesional, pero VOS validás. Un test plan generado que no
  leíste no es un test plan: es un riesgo con formato lindo.
- **No garantiza corrección.** Los LLMs se equivocan con confianza.
  Parte del curso es aprender a detectar cuándo.
- **No tiene las capas enterprise de QATES**: enforcement por hooks,
  métricas de adopción, supervivencia de tests, regresión de prompts.
  Acá está el motor; el chasis de producción lo ves en QATES.
- **No corre solo.** Cada acción pasa por el sistema de permisos de
  Claude Code: el agente propone, vos aprobás.
- **No escribe en Jira, ADO ni Confluence.** QASON puede LEER un ticket
  que le pegues, pero no comenta, no transiciona y no publica páginas.
  Es a propósito: aprendés en un entorno donde no podés romper el
  tracker de tu empresa sin querer. Esas capacidades viven en QATES.

## 2. Requisitos

1. **Claude Code** instalado y autenticado ([claude.com/claude-code](https://claude.com/claude-code)).
2. **Go 1.26+** para compilar el instalador ([go.dev/dl](https://go.dev/dl)).
   Instalalo **antes** de la clase: con una versión más vieja Go descarga el
   toolchain correcto solo, pero son ~80 MB y el wifi del aula no perdona.
3. Un proyecto donde practicar (cualquier repo con código sirve; los
   ejercicios incluyen tickets de ejemplo).

## 3. Instalación

```bash
git clone https://github.com/FacundoPasqua/Qason.git
cd Qason
go run ./cmd/qason        # wizard interactivo
```

Elegí **Install** y listo. Si preferís sin wizard: `go run ./cmd/qason install`.

**Verificá la instalación**: abrí Claude Code en cualquier proyecto y
escribí "¿qué agentes QA tenés disponibles?" — debería nombrar
`qa-analyst`, `qa-test-designer` y `qa-automator`.

**Para desinstalar**: `go run ./cmd/qason uninstall` (tu contenido
propio de `CLAUDE.md` queda intacto).

### Si usás GitHub Copilot

VS Code lee sub-agentes de `~/.claude/agents` y skills de
`~/.claude/skills`, así que el install de arriba ya te dejó todo. Lo
único que Copilot no lee es `CLAUDE.md` — donde vive el orquestador.
Corré esto en la raíz de tu proyecto:

```bash
go run ./cmd/qason install --copilot .
```

Y pensá por qué hizo falta, porque es la lección escondida acá: **un
agente no es un modelo, es un modelo MÁS su configuración**. Cambiás la
herramienta y los prompts viajan igual; lo que no viaja es dónde cada
herramienta va a buscarlos. Por eso importa entender la arquitectura y
no solo copiar comandos.

Dos diferencias: los colores de los agentes probablemente no se vean
(Copilot no soporta ese campo) y `model: sonnet` es nomenclatura de
Claude Code. Probalo antes de la clase.

### Qué instaló, exactamente

| Qué | Dónde | Para qué |
|---|---|---|
| 3 sub-agentes | `~/.claude/agents/qa-*.md` | Los especialistas que Claude Code puede invocar |
| 31 skills | `~/.claude/skills/qason/` | Capacidades atómicas que los agentes cargan según la tarea |
| Orquestador | bloque en `~/.claude/CLAUDE.md` | Las reglas de ruteo y síntesis |

Todo es **markdown legible**. Abrí los archivos. En serio: el prompt ES
el código de un agente, y leerlo es la mejor clase de ingeniería de
prompts que vas a tener.

## 4. Los tres agentes, en detalle

### `qa-analyst` — el que piensa antes de testear

Transforma requisitos ambiguos en estrategia: requisitos explícitos e
implícitos, criterios faltantes, preguntas para el PO, test plan y
matriz de riesgo. Su regla de oro: **nunca asumas que los requisitos
están completos**.

Probalo: *"Analizá este ticket: [pegá un ticket]"*

### `qa-test-designer` — el que convierte estrategia en casos

Toma el test plan y produce casos ejecutables: Given/When/Then,
prioridades según riesgo, datos de prueba, casos negativos y de borde.
Su regla de oro: **cada caso prueba UNA sola cosa**.

Probalo: *"Diseñá casos de prueba para el flujo de login con lockout"*

### `qa-automator` — el que escribe y CORRE el código

Convierte casos en tests automatizados. Detecta el framework del
proyecto, respeta sus convenciones, y — esto es clave — **ejecuta lo
que genera**. Un test que nunca corrió no es un test. Si falla, intenta
arreglarlo con un límite de 3 intentos y después te lo entrega con un
handoff estructurado de lo que probó.

Probalo: *"Generá tests unitarios para [archivo]"*

## 5. El pipeline Spec-to-Test (la clase completa en un comando)

Pedile a Claude Code:

> Analizá este ticket y creá los tests: [ticket]

Y mirá la secuencia — **literalmente miralla**: cada agente tiene su
color en la task list de Claude Code, así que el pipeline se lee de un
vistazo.

| | Agente | Qué está haciendo cuando lo ves |
|---|---|---|
| 🔵 cyan | `qa-analyst` | Está pensando: requisitos, riesgos, preguntas |
| 🟢 green | `qa-test-designer` | Está diseñando: casos, datos, prioridades |
| 🟣 purple | `qa-automator` | Está ejecutando: código y tests corriendo |

Si ves violeta antes que cyan, algo se salteó el pipeline — y esa es
una observación que vale oro cuando debuguees tus propios agentes.

El orquestador delega al analyst → el análisis
alimenta al designer → los casos alimentan al automator → síntesis
final con los resultados de los tres. Cada agente recibe el output del
anterior como contexto. Eso es **orquestación agéntica**, y acabás de
verla funcionar en tu terminal.

## 6. Ejercicios

### Ejercicio 1 — El ticket con trampa (análisis)

Pegale este ticket al agente y pedile que lo analice:

> Como usuario quiero recuperar mi contraseña por email.
> Criterios: (1) link "olvidé mi contraseña" en el login, (2) se envía
> un mail con link de reseteo, (3) si el email no existe se muestra
> "No existe cuenta para este email".

**Antes de leer la respuesta**, anotá vos: ¿qué está mal en este
ticket? Hay un problema de seguridad escrito DENTRO de los criterios de
aceptación. Después compará con lo que encontró el `qa-analyst`.
¿Lo detectó? ¿Detectó cosas que vos no? ¿Vos viste algo que él no?

### Ejercicio 2 — Del análisis a los casos (diseño)

Con el análisis del ejercicio 1: *"Diseñá los casos de prueba
priorizados por la matriz de riesgo"*. Evaluá el output: ¿los casos de
seguridad quedaron primeros? ¿Hay casos negativos? ¿Cada caso prueba
una sola cosa?

### Ejercicio 3 — Tests que corren (automatización)

En un proyecto real (o uno de práctica): *"Generá tests unitarios para
[una función con lógica]"*. Observá: ¿detectó el framework? ¿Los
ejecutó? ¿Qué hizo cuando fallaron?

### Ejercicio 4 — El pipeline completo

Tomá un ticket real de tu trabajo (anonimizalo), corré el pipeline
completo y auditá el resultado de punta a punta como si fueras el QA
lead que revisa el trabajo de tres juniors brillantes pero nuevos.
Porque eso es exactamente lo que son.

### Ejercicio 5 — Leer los prompts (el meta-ejercicio)

Abrí `~/.claude/agents/qa-analyst.md` y `~/.claude/skills/qason/prd-analyzer/SKILL.md`.
Preguntate: ¿por qué el prompt insiste en requisitos IMPLÍCITOS? ¿por
qué el formato de salida es una tabla con esas columnas? Cada regla de
esos archivos existe porque un LLM sin ella hace las cosas mal. Eso es
ingeniería de prompts: reglas nacidas de fallas observadas.

## 7. Problemas comunes

| Síntoma | Causa probable | Solución |
|---|---|---|
| Claude Code no ve los agentes | Instalaste con Claude Code abierto | Reiniciá Claude Code (lee `agents/` al arrancar) |
| "command not found: go" | Go no instalado / no en PATH | Instalá desde go.dev/dl y reabrí la terminal |
| El agente responde genérico, sin delegar | El bloque QASON no está en `CLAUDE.md` | Corré `go run ./cmd/qason install` de nuevo (es idempotente) |
| Los tests generados no corren | Falta el framework en el proyecto | El automator te lo va a decir — leé su handoff |

## 8. La regla de la casa

La IA es una herramienta. **Nosotros dirigimos, la IA ejecuta.** Todo
lo que un agente produzca pasa por tu criterio antes de llegar a un
repo, un Jira o un cliente. Si no entendés lo que generó, no lo uses:
primero entendé, después usá. Los atajos en el aprendizaje se pagan
con intereses.

¿Y cuando QASON te quede chico? El siguiente paso es
[QATES](https://github.com/FacundoPasqua/Qates): 5 agentes, 51 skills,
workflows compuestos, enforcement mecánico y métricas de adopción. Todo
lo que aprendiste acá aplica directo.
