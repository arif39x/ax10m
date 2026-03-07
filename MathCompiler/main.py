import re
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import sympy as sp

app = FastAPI()

class EquationRequest(BaseModel):
    equation: str

x, y, z = sp.symbols('x y z', real=True)
state_x, state_y, state_z = sp.symbols('state_x state_y state_z', real=True)

ALLOWED_VARS = {"x": x, "y": y, "z": z, "state.x": state_x, "state.y": state_y, "state.z": state_z}

def to_wgsl(expr) -> str:
    c_code = sp.ccode(expr)
    c_code = re.sub(r'(?<![a-zA-Z0-9_\.])(\d+)(?![a-zA-Z0-9_\.])', r'\1.0', c_code)
    c_code = c_code.replace('state_x', 'state.x')
    c_code = c_code.replace('state_y', 'state.y')
    c_code = c_code.replace('state_z', 'state.z')

    return f"return {c_code};"

@app.post("/compile_sdf")
async def compile_sdf(req: EquationRequest):
    try:
        eq_str = req.equation.replace('state.x', 'state_x').replace('state.y', 'state_y').replace('state.z', 'state_z')
        expr = sp.sympify(eq_str, locals=ALLOWED_VARS)
        
        free_symbols = expr.free_symbols
        allowed_set = set(ALLOWED_VARS.values())
        if not free_symbols.issubset(allowed_set):
            unauthorized = free_symbols - allowed_set
            raise HTTPException(status_code=400, detail=f"Unauthorized variables used: {unauthorized}")
            
        wgsl_code = to_wgsl(expr)
        
        return {"status": "success", "wgsl": wgsl_code}
        
    except Exception as e:
        raise HTTPException(status_code=400, detail=f"Compilation error: {str(e)}")
