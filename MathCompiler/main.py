import re
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import sympy as sp
from ast_firewall import validate_equation

app = FastAPI()

class EquationRequest(BaseModel):
    equation: str

x, y, z = sp.symbols('x y z', real=True)
state_x, state_y, state_z = sp.symbols('state_x state_y state_z', real=True)

ALLOWED_VARS = {"x": x, "y": y, "z": z, "state_x": state_x, "state_y": state_y, "state_z": state_z}

def to_wgsl_module(expr) -> str:
    c_code = sp.ccode(expr)
    c_code = re.sub(r'(?<![a-zA-Z0-9_\.])(\d+)(?![a-zA-Z0-9_\.])', r'\1.0', c_code)
    c_code = c_code.replace('state_x', 'state.x')
    c_code = c_code.replace('state_y', 'state.y')
    c_code = c_code.replace('state_z', 'state.z')

    return f"""
struct State {{
    x: f32,
    y: f32,
    z: f32,
    padding: f32,
}}

struct Emitter {{
    x: f32,
    y: f32,
    z: f32,
    amplitude: f32,
    sigma: f32,
}}

@group(0) @binding(0)
var<uniform> state: State;

struct VertexOutput {{
    @builtin(position) clip_position: vec4<f32>,
    @location(0) uv: vec2<f32>,
}};

@vertex
fn vs_main(@builtin(vertex_index) in_vertex_index: u32) -> VertexOutput {{
    var out: VertexOutput;
    let x = f32((in_vertex_index << 1) & 2u);
    let y = f32(in_vertex_index & 2u);
    out.clip_position = vec4<f32>(x * 2.0 - 1.0, 1.0 - y * 2.0, 0.0, 1.0);
    out.uv = vec2<f32>(x, y);
    return out;
}}

fn map(p: vec3<f32>) -> vec2<f32> {{
    let x = p.x;
    let y = p.y;
    let z = p.z;
    let dist = {c_code};
    return vec2<f32>(dist, 1.0); // 1.0 = Default Material
}}

fn get_force_intensity(p: vec3<f32>) -> f32 {{
    let emitter_pos = vec3<f32>(0.0, 0.0, -20.0);
    let dist = length(p - emitter_pos);
    let sigma = 50.0;
    return exp(-pow(dist, 2.0) / (2.0 * pow(sigma, 2.0)));
}}

@fragment
fn fs_main(in: VertexOutput) -> @location(0) vec4<f32> {{
    let uv = in.uv * 2.0 - 1.0;

    let ro = vec3<f32>(0.0, 0.0, 200.0);
    let rd = normalize(vec3<f32>(uv.x, uv.y, -1.0));

    var t = 0.0;
    var force_acc = 0.0;

    for (var i = 0; i < 64; i = i + 1) {{
        let p = ro + rd * t;
        let res = map(p);
        let d = res.x;
        let mat = res.y;

        force_acc += get_force_intensity(p) * 0.05;

        if (d < 0.001) {{
            var col = vec3<f32>(0.5, 0.8, 1.0);
            if (mat > 1.5) {{ // Example for other materials
                col = vec3<f32>(1.0, 0.5, 0.2);
            }}

            col = col * (1.0 - f32(i)/64.0);
            let force_glow = vec3<f32>(0.2, 0.4, 1.0) * force_acc;
            return vec4<f32>(col + force_glow, 1.0);
        }}
        t = t + d;
        if (t > 400.0) {{
            break;
        }}
    }}

    let field_col = vec3<f32>(0.0, 0.05, 0.1) * force_acc;
    return vec4<f32>(field_col, 1.0);
}}
"""
@app.post("/compile_sdf")
async def compile_sdf(req: EquationRequest):
    try:
        validate_equation(req.equation)
        eq_str = req.equation.replace('state.x', 'state_x').replace('state.y', 'state_y').replace('state.z', 'state_z')
        expr = sp.sympify(eq_str, locals=ALLOWED_VARS)

        free_symbols = expr.free_symbols
        allowed_set = set(ALLOWED_VARS.values())
        if not free_symbols.issubset(allowed_set):
            unauthorized = free_symbols - allowed_set
            raise HTTPException(status_code=400, detail=f"Unauthorized variables used: {unauthorized}")

        wgsl_code = to_wgsl_module(expr)

        return {"status": "success", "wgsl": wgsl_code}

        
    except Exception as e:
        raise HTTPException(status_code=400, detail=f"Compilation error: {str(e)}")
