package cr.aoc2023.day13;

import java.util.List;

public class Pattern {
    private int[][] data;

/**
* Takes a list of Strings and creates the internal data array
* @param input lines of the format "#.##....##.#.
*/
    public Pattern(List<String> input) {
        data = new int[input.size()][input.get(0).length()];
        int row = 0;
        for (String line: input) {
            int col = 0;
            for (byte b: line.getBytes()) {
                switch (b) {
                    case '#': data[row][col] = 1; break;
                    case '.': data[row][col] = 0; break;
                    default: System.out.println("Error parsing: " + line);
                }
                col++;
            }
            row++;
        }
    }

    public int getRows() {
        return data.length;
    }

    public int getCols() {
        return data[0].length;
    }

    public void transposeMatrix() {
        int m = data.length;
        int n = data[0].length;
        int[][] transposedMatrix = new int[n][m];
        for(int x=0; x < n; x++) {
            for(int y = 0; y < m; y++) {
                transposedMatrix[x][y]= data[y][x];
            }
        }
        data = transposedMatrix;
    }

    public int[] getRowPlusRows(int row, int delta) { 
        int[] result = new int[getCols()*delta];
        int idx=0;
        for (int r=0; r<delta; r++) {
            for (int c=0; c<getCols(); c++)
                result[idx++] = data[row+r][c];
        }
    
       return result;
    }


    public int[] getRowMinusRows(int row, int delta) { 
        int[] result = new int[getCols()*delta];
        int idx=0;
        for (int r=0; r<delta; r++) {
            for (int c=0; c<getCols(); c++)
                result[idx++] = data[row-r][c];
        }
        return result;
    }
}
